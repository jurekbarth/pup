package project

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/jurekbarth/pup/worker"
)

// Deploy informations
type Deploy struct {
	ID        string    `json:"id"`
	Stage     string    `json:"stage"`
	Timestamp time.Time `json:"timestamp"`
}

// User for cognito
type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Groups   []string `json:"groups"`
}

// Rule for lambda edge server
type Rule struct {
	Groups []string `json:"group-permissions"`
}

// Project config
type Project struct {
	Version      int               `json:"version"`
	CustomerCode string            `json:"customer-code"`
	Project      string            `json:"project"`
	Root         string            `json:"root"`
	Users        []User            `json:"users"`
	Rules        []map[string]Rule `json:"rules"`
	AWSProfile   string            `json:"aws-profile"`
	AWSBucket    string            `json:"aws-bucket"`
	AWSRegion    string            `json:"aws-region"`
}

// Read reads the downloaded project config
func Read(w *worker.Worker) (*Project, *Deploy, error) {
	dir := w.Config.UnzipDir
	return read(dir)
}

func read(dir string) (*Project, *Deploy, error) {
	raw, err := ioutil.ReadFile(dir + "/pup.json")
	if err != nil {
		return nil, nil, err
	}
	var project Project
	err = json.Unmarshal(raw, &project)
	if err != nil {
		return nil, nil, err
	}

	raw, err = ioutil.ReadFile(dir + "/deployConfig.json")
	if err != nil {
		return nil, nil, err
	}
	var projectDeploy Deploy
	err = json.Unmarshal(raw, &projectDeploy)
	if err != nil {
		return nil, nil, err
	}
	return &project, &projectDeploy, nil
}
