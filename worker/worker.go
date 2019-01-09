package worker

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jurekbarth/pup/worker/config"
	"github.com/jurekbarth/pup/worker/event"
)

// Config for a project.
type Config = config.Config

// ReadConfig reads the configuration from `path`.
var ReadConfig = config.ReadConfig

// GetSQSMessage parses sqs message
// var GetSQSMessage = sqs.GetSQSMessage

// Worker ...
type Worker struct {
	Config *Config
	Events event.Events
}

// New ...
func New(c *Config, events event.Events) *Worker {
	return &Worker{
		Config: c,
		Events: events,
	}
}

// MakeSession returns a aws session with a region
func MakeSession(w *Worker, region *string) (*session.Session, error) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	r := region
	if r == nil {
		r = &w.Config.AWSRegion
	}
	// Specify profile to load for the session's config
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: r},
		// Profile: w.Config.AWSProfile,
	})
	return sess, err
}
