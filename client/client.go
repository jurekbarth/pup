package client

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"


	"github.com/jurekbarth/pup/client/config"
	"github.com/jurekbarth/pup/client/event"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

// Config for a project.
type Config = config.Config

// ReadConfig reads the configuration from `path`.
var ReadConfig = config.ReadConfig

type Project struct {
	config  *Config
	events  event.Events
	zipFile *string
}

type Deploy struct {
	ID        string    `json:"id"`
	Stage     string    `json:"stage"`
	Timestamp time.Time `json:"timestamp"`
}

// New ...
func New(c *Config, events event.Events) *Project {
	return &Project{
		config: c,
		events: events,
	}
}

// Deploy the project.
func (p *Project) Deploy(d Deploy) error {
	if err := p.zip(d); err != nil {
		return errors.Wrap(err, "zipping")
	}
	if err := p.upload(d); err != nil {
		return errors.Wrap(err, "uploading")
	}
	cloudfrontID, cacheID, err := p.logs(d)
	if err != nil {
		return errors.Wrap(err, "worker")
	}
	if cloudfrontID != "" {
		err = p.waitCloudfront(cloudfrontID)
		if err != nil {
			return errors.Wrap(err, "cloudfront")
		}
	}
	if cacheID != "" {
		err = p.waitInvalidation(cacheID)
		if err != nil {
			return errors.Wrap(err, "cloudfront invalidation")
		}
	}
	duration := time.Second * 2
	time.Sleep(duration)

	return nil
}

// cache id has the format of "I3J6MXXQ4Y2SN6###E3CP9NCBEQTOMZ" (cachid###distributionid)
func (p *Project) waitInvalidation(cacheID string) error {
	cid := cacheID
	l := len([]rune(cid))
	idx := strings.Index(cid, "###")
	distributionID := cid[idx+3 : l]
	invalidationID := cid[0:idx]
	e := event.Event{
		Name: "cfinvalidation.start",
	}
	p.events.Emit(e)
	if p.config.AWSProfile != "" {
		setProfile(p.config.AWSProfile)
	}
	s := session.New(aws.NewConfig().WithRegion(p.config.AWSRegion))
	cfs := cloudfront.New(s)
	input := &cloudfront.GetInvalidationInput{
		DistributionId: aws.String(distributionID),
		Id:             aws.String(invalidationID),
	}
retry:
	invalidation, err := cfs.GetInvalidation(input)
	if err != nil {
		return err
	}
	if invalidation.Invalidation.Status == nil {
		return errors.New("No invalidation state found")
	}
	if *invalidation.Invalidation.Status == "Completed" {
		e := event.Event{
			Name: "cfinvalidation.done",
		}
		p.events.Emit(e)
		return nil
	}
	duration := time.Second * 15
	time.Sleep(duration)
	goto retry
}

func (p *Project) waitCloudfront(cloudfrontID string) error {
	e := event.Event{
		Name: "cf.start",
	}
	p.events.Emit(e)
	if p.config.AWSProfile != "" {
		setProfile(p.config.AWSProfile)
	}
	s := session.New(aws.NewConfig().WithRegion(p.config.AWSRegion))
	cfs := cloudfront.New(s)
	input := &cloudfront.GetDistributionInput{
		Id: aws.String(cloudfrontID),
	}
retry:
	cf, err := cfs.GetDistribution(input)
	if err != nil {
		return err
	}
	if cf.Distribution.Status == nil {
		return errors.New("cannot get status of cloudfront distribution")
	}
	if cf.Distribution.Status != nil && *cf.Distribution.Status == "Deployed" {
		e := event.Event{
			Name: "cf.done",
		}
		p.events.Emit(e)
		return nil
	}
	duration := time.Second * 30
	time.Sleep(duration)
	goto retry
}

func (p *Project) logs(d Deploy) (cloudfrontID string, cacheID string, err error) {
	e := event.Event{
		Name: "logs.start",
	}
	p.events.Emit(e)
	// filename := "1529940836-vi-biotope-frontend-latestt.zip"
	filename := *p.zipFile
	if p.config.AWSProfile != "" {
		setProfile(p.config.AWSProfile)
	}

	s := session.New(aws.NewConfig().WithRegion(p.config.AWSRegion))
	client := dynamodb.New(s)

	l := len([]rune(filename))
	id := filename[0 : l-4]
	length := 0
	first := true
	invalidationID := ""
	distributionID := ""
retry:
	logs, err := getLogs(p, id, client)
	if err != nil {
		return "", "", err
	}
	if len(logs.Logs) != length {
		if first {
			e := event.Event{
				Name: "logs.done",
			}
			p.events.Emit(e)
		}
		first = false
		for idx, log := range logs.Logs {
			if idx >= length {
				if log.Error != "" {
					return "", "", errors.New(log.Error)
				}
				e := event.Event{
					Name:  log.Name,
					Value: log.Value,
				}
				p.events.Emit(e)
				if log.Name == "cloudfrontupdate.done" {
					distributionID = log.Value[14:]
				}
				if log.Name == "cloudfront.done" {
					invalidationID = log.Value[15:]
				}
				if log.Name == "sqs.done" {
					return distributionID, invalidationID, nil
				}
			}
		}
		length = len(logs.Logs)
	}
	duration := time.Second
	time.Sleep(duration)
	goto retry
}

type Logs struct {
	DeployID string `json:"deployid"`
	Logs     []Log
}

type Log struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Error string `json:"error"`
}

func getLogs(p *Project, id string, client *dynamodb.DynamoDB) (*Logs, error) {

	key := map[string]*dynamodb.AttributeValue{
		"deployid": {
			S: aws.String(id),
		},
	}
	res, err := client.GetItem(&dynamodb.GetItemInput{
		TableName:      aws.String(p.config.AWSDynamoDBLogsTable),
		Key:            key,
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	logs := Logs{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &logs)
	if err != nil {
		return nil, err
	}
	return &logs, nil
}

func (p *Project) upload(d Deploy) error {
	e := event.Event{
		Name: "upload.start",
	}
	p.events.Emit(e)
	if p.config.AWSProfile != "" {
		setProfile(p.config.AWSProfile)
	}

	s := session.New(aws.NewConfig().WithRegion(p.config.AWSRegion))
	uploader := s3manager.NewUploaderWithClient(s3.New(s))

	bucket := aws.String(p.config.AWSBucket)
	filename := p.zipFile

	file, err := os.Open(*filename)
	if err != nil {
		return errors.Wrap(err, "uploading file")
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    filename,
		Body:   fileBytes,
	})

	if err != nil {
		return errors.Wrap(err, "uploading file")
	}
	e = event.Event{
		Name: "upload.done",
	}
	p.events.Emit(e)
	return nil
}

func setProfile(name string) {
	os.Setenv("AWS_PROFILE", name)
}

func (p *Project) zip(d Deploy) error {
	e := event.Event{
		Name: "zip.start",
	}
	p.events.Emit(e)
	ts := strconv.FormatInt(d.Timestamp.Unix(), 10)
	p.zipFile = aws.String(ts + "-" + p.config.CustomerCode + "-" + p.config.Project + "-" + d.Stage + ".zip")
	err := Zip(p.config.Root, "./"+*p.zipFile, d)
	e = event.Event{
		Name: "zip.done",
	}
	p.events.Emit(e)
	return errors.Wrap(err, "zipping file")
}

// Zip the file
func Zip(dir string, dest string, d Deploy) error {
	baseFolder := dir
	if filepath.Ext(baseFolder) == "" {
		baseFolder = filepath.Clean(baseFolder) + "/"
	}

	// Get a Buffer to Write To
	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	err = addFiles(w, baseFolder, "")

	if err != nil {
		return err
	}

	err = addFile(w, "./", "/", "pup.json")
	if err != nil {
		return err
	}

	l := len([]rune(dest))
	d.ID = dest[2 : l-4]
	// create file
	cd, err2 := json.Marshal(&d)
	if err2 != nil {
		return err2
	}
	err = addData(w, "./", "/", "deployConfig.json", cd)
	if err != nil {
		return err
	}
	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func addFiles(w *zip.Writer, basePath string, baseInZip string) error {

	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			err = addFile(w, basePath, baseInZip, file.Name())
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			bs := baseInZip + file.Name() + "/"
			addFiles(w, newBase, bs)
		}
	}
	return nil
}

func addFile(w *zip.Writer, basePath string, baseInZip string, fileName string) error {
	data, err := ioutil.ReadFile(basePath + fileName)
	if err != nil {
		return err
	}
	return addData(w, basePath, baseInZip, fileName, data)
}

func addData(w *zip.Writer, basePath string, baseInZip string, fileName string, data []byte) error {
	// Add some files to the archive.
	f, err := w.Create(baseInZip + fileName)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
