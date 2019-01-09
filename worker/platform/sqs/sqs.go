package sqs

import (
	"encoding/json"
	"errors"

	"github.com/jurekbarth/pup/worker"
	"github.com/jurekbarth/pup/worker/platform/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// GetSqsMessage SQS Thingy
func GetSqsMessage(w *worker.Worker) (*sqs.ReceiveMessageOutput, error) {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return nil, err
	}
	svc := sqs.New(session)
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(w.Config.AWSqsURI),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(900), // 15 minutes
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		return nil, err
	}
	if len(result.Messages) == 0 {
		return nil, errors.New("Received no messages")
	}
	return result, nil
}

// DeleteSqsMessage deletes a SQS Message
func DeleteSqsMessage(w *worker.Worker, result *sqs.ReceiveMessageOutput) (*sqs.DeleteMessageOutput, error) {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return nil, err
	}
	svc := sqs.New(session)
	resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(w.Config.AWSqsURI),
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})

	if err != nil {
		return nil, err
	}
	return resultDelete, nil
}

// GetS3Key returns a S3 Key for and S3 Event in SQS
func GetS3Key(result *sqs.ReceiveMessageOutput) (string, error) {
	s3Message := new(s3.S3Event)
	if err := json.Unmarshal([]byte(*result.Messages[0].Body), &s3Message); err != nil {
		return "", err
	}
	return s3Message.Records[0].S3.Object.Key, nil
}
