package s3

import (
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jurekbarth/pup/worker"
)

// Download ...
func Download(w *worker.Worker, key string) error {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return err
	}
	filename := filepath.Base(key)
	f, err := os.Create(w.Config.DownloadDir + "/" + filename)
	if err != nil {
		return err
	}
	s3Svc := s3manager.NewDownloader(session)
	params := &s3.GetObjectInput{
		Bucket: aws.String(w.Config.AWSFromBucket),
		Key:    aws.String(key),
	}
	_, err = s3Svc.Download(f, params)
	if err != nil {
		return err
	}
	return nil
}

// Upload ...
func Upload(w *worker.Worker, projectPath string) error {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return err
	}
	absUnzipDir, err := filepath.Abs(w.Config.UnzipDir)
	if err != nil {
		return err
	}
	if isDirectory(w.Config.UnzipDir) {
		err := uploadDirToS3(absUnzipDir, session, w.Config.AWSDestinationBucket, projectPath)
		if err != nil {
			return err
		}
	} else {
		relativeFilePath := filepath.Base(absUnzipDir)
		err := UploadFileToS3(absUnzipDir, session, w.Config.AWSDestinationBucket, projectPath, relativeFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func uploadDirToS3(absUnzipDir string, session *session.Session, AWSDestinationBucket string, projectPath string) error {
	fileList := []string{}
	filepath.Walk(absUnzipDir, func(path string, f os.FileInfo, err error) error {
		if isDirectory(path) {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})

	// base := filepath.Base(dirPath)
	for _, filePath := range fileList {
		relativeFilePath, err := filepath.Rel(absUnzipDir, filePath)
		if err != nil {
			return err
		}
		err = UploadFileToS3(filePath, session, AWSDestinationBucket, projectPath, relativeFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadFileToS3 ...
func UploadFileToS3(filePath string, session *session.Session, AWSDestinationBucket string, projectPath string, relativeFilePath string) error {
	// An s3 service
	s3Svc := s3.New(session)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	var key string
	key = projectPath + "/" + relativeFilePath
	// Upload the file to the s3 given bucket
	fileExtension := filepath.Ext(relativeFilePath)
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".ico", "image/x-icon")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".mpeg", "video/mpeg")
	mime.AddExtensionType(".zip", "application/zip")
	mime.AddExtensionType(".webm", "video/webm")
	mime.AddExtensionType(".webp", "video/webp")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".woff2", "font/woff2")
	typ := mime.TypeByExtension(fileExtension)
	params := &s3.PutObjectInput{
		Bucket:      aws.String(AWSDestinationBucket), // Required
		Key:         aws.String(key),                  // Required
		Body:        file,
		ContentType: aws.String(typ),
	}
	_, err = s3Svc.PutObject(params)
	if err != nil {
		return err
	}
	return nil
}

func isDirectory(path string) bool {
	fd, err := os.Stat(path)
	if err != nil {
		log.Panicln(err)
		os.Exit(2)
	}
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}

// DeleteDir ....
func DeleteDir(w *worker.Worker, dirPath string) error {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return err
	}
	svc := s3.New(session)
delete:
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(w.Config.AWSDestinationBucket),
		Prefix:  aws.String(dirPath),
		MaxKeys: aws.Int64(200),
	}

	// Create a delete list objects iterator
	iter := s3manager.NewDeleteListIterator(svc, input)
	// Create the BatchDelete client
	batcher := s3manager.NewBatchDeleteWithClient(svc)

	if err := batcher.Delete(aws.BackgroundContext(), iter); err != nil {
		return err
	}
	result, err2 := svc.ListObjects(input)
	if err2 != nil {
		return err2
	}
	if len(result.Contents) > 0 {
		goto delete
	}
	return nil
}
