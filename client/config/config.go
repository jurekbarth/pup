package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jurekbarth/pup/client/internal/validate"
	"github.com/pkg/errors"
)

type user struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Groups   []string `json:"groups"`
}

// Rule for lambda edge server
type Rule struct {
	Groups []string `json:"group-permissions"`
}

// Config for the project.
type Config struct {
	Version              int               `json:"version"`
	CustomerCode         string            `json:"customer-code"`
	Project              string            `json:"project"`
	Root                 string            `json:"root"`
	Users                []user            `json:"users"`
	Rules                []map[string]Rule `json:"rules"`
	AWSProfile           string            `json:"aws-profile"`
	AWSBucket            string            `json:"aws-bucket"`
	AWSRegion            string            `json:"aws-region"`
	AWSDynamoDBLogsTable string            `json:"aws-dynamodb-logs-table"`
}

type validator interface {
	Validate() error
}

// regquiring rules is important as there's no default config for that
func requiredRules(m []map[string]Rule) error {
	for _, ruleWrapper := range m {
		for _, r := range ruleWrapper {
			if err := validate.RequiredStrings(r.Groups); err != nil {
				return err
			}
		}
	}
	return nil
}

// Validate implementation.
func (c *Config) Validate() error {
	if err := validate.RequiredString(c.Project); err != nil {
		return errors.Wrap(err, ".project")
	}

	if err := validate.Name(c.Project); err != nil {
		return errors.Wrapf(err, ".name %q", c.Project)
	}

	if err := requiredRules(c.Rules); err != nil {
		return errors.Wrap(err, ".rules")
	}

	if err := validate.RequiredString(c.CustomerCode); err != nil {
		return errors.Wrap(err, ".customer-code")
	}

	if err := validate.Name(c.CustomerCode); err != nil {
		return errors.Wrapf(err, ".customer-code %q", c.CustomerCode)
	}

	return nil
}

// Default Config
func (c *Config) Default() error {

	if c.Root == "" {
		c.Root = "./dist"
	}

	if c.AWSProfile == "" {
		c.AWSProfile = "biotope"
	}

	if c.AWSBucket == "" {
		c.AWSBucket = "biotope-zip-bucket"
	}

	if c.AWSRegion == "" {
		c.AWSRegion = "eu-central-1"
	}

	if c.AWSDynamoDBLogsTable == "" {
		c.AWSDynamoDBLogsTable = "biotope-logs-table"
	}

	return nil
}

// ParseConfig returns config from JSON bytes.
func ParseConfig(b []byte) (*Config, error) {
	c := &Config{}

	if err := json.Unmarshal(b, c); err != nil {
		return nil, errors.Wrap(err, "parsing json")
	}

	if err := c.Default(); err != nil {
		return nil, errors.Wrap(err, "defaulting")
	}

	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating")
	}

	return c, nil
}

// ReadConfig reads the configuration from `path`.
func ReadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(b)
}
