package config

import (
	"fmt"

	"github.com/jurekbarth/getenv"
	"github.com/jurekbarth/pup/worker/internal/validate"
	"github.com/pkg/errors"
)

// Config for pupserver commandline
type Config struct {
	DownloadDir                     string `envconfig:"DOWNLOAD_DIR"`
	UnzipDir                        string `envconfig:"UNZIP_DIR"`
	AWSFromBucket                   string `envconfig:"S3_FROM_BUCKET"`
	AWSDestinationBucket            string `envconfig:"S3_DESTINATION_BUCKET"`
	AWSRegion                       string `envconfig:"DEFAULT_REGION"`
	AWSqsURI                        string `envconfig:"SQS_URI"`
	AWSCloudfrontID                 string `envconfig:"CF_ID"`
	AWSCloudfrontDomain             string `envconfig:"CF_DOMAIN"`
	AWSDynamoDBRulesTable           string `envconfig:"DDB_RULES_TABLE"`
	AWSDynamoDBLogsTable            string `envconfig:"DDB_LOGS_TABLE"`
	AWSLambdaBucket                 string `envconfig:"S3_LAMBDA_BUCKET"`
	AWSLambdaARN                    string `envconfig:"LAMBDA_EDGE_ARN"`
	AWSCognitoClientBackendClientID string `envconfig:"COGNITO_CLIENT_BACKEND_CLIENT_ID"`
	AWSCognitoClientSPAClientID     string `envconfig:"COGNITO_CLIENT_SPA_CLIENT_ID"`
	AWSCognitoClientSPASubDomain    string `envconfig:"COGNITO_CLIENT_SPA_SUBDOMAIN"`
	AWSCognitoPoolID                string `envconfig:"COGNITO_POOL_ID"`
	EmailDomain                     string `envconfig:"EMAIL_DOMAIN"`
}

// Validate implementation.
func (c *Config) Validate() error {
	if err := validate.RequiredString(c.DownloadDir); err != nil {
		return errors.Wrap(err, ".DownloadDir")
	}

	if err := validate.RequiredString(c.UnzipDir); err != nil {
		return errors.Wrap(err, ".UnzipDir")
	}

	if err := validate.RequiredString(c.AWSFromBucket); err != nil {
		return errors.Wrap(err, ".AWSFromBucket")
	}

	if err := validate.RequiredString(c.AWSDestinationBucket); err != nil {
		return errors.Wrap(err, ".AWSDestinationBucket")
	}

	if err := validate.RequiredString(c.AWSqsURI); err != nil {
		return errors.Wrap(err, ".AWSqsURI")
	}

	if err := validate.RequiredString(c.AWSCloudfrontID); err != nil {
		return errors.Wrap(err, ".AWSCloudfrontID")
	}

	if err := validate.RequiredString(c.AWSDynamoDBRulesTable); err != nil {
		return errors.Wrap(err, ".AWSDynamoDBTable")
	}

	if err := validate.RequiredString(c.AWSDynamoDBLogsTable); err != nil {
		return errors.Wrap(err, ".AWSDynamoDBLogsTable")
	}

	if err := validate.RequiredString(c.AWSLambdaBucket); err != nil {
		return errors.Wrap(err, ".AWSLambdaBucket")
	}

	if err := validate.RequiredString(c.AWSLambdaARN); err != nil {
		return errors.Wrap(err, ".AWSLambdaARN")
	}

	return nil
}

// ReadConfig reads the configuration from `path`.
func ReadConfig() (*Config, error) {
	var c Config
	getenv.Process("PUP_", &c)
	fmt.Println(c)
	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating")
	}

	return &c, nil
}
