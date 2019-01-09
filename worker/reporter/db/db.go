package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jurekbarth/pup/worker"

	"github.com/jurekbarth/pup/worker/event"
)

// Report upates a db entry for client <-- worker events
func Report(events <-chan *event.Event, w *worker.Worker) {
	r := reporter{
		worker:  w,
		events:  events,
		reports: []report{},
	}

	r.Start()
}

type reporter struct {
	worker  *worker.Worker
	events  <-chan *event.Event
	reports []report
}

type report struct {
	timestamp time.Time
	name      string
	errorMsg  string
}

// Start handling events.
func (r *reporter) Start() {
	session, err := worker.MakeSession(r.worker, nil)
	if err != nil {
		fmt.Println(err)
	}
	client := ddb.New(session)
	tableName := r.worker.Config.AWSDynamoDBLogsTable
	for {
		select {
		case e := <-r.events:
			fmt.Println(e.ID + " " + e.Name + ": " + e.Value)
			logToDB(client, e, tableName)
		}
	}
}

func logToDB(client *ddb.DynamoDB, e *event.Event, tableName string) {
	key := map[string]*ddb.AttributeValue{
		"deployid": {
			S: aws.String(e.ID),
		},
	}
	logObject := map[string]*ddb.AttributeValue{
		"name": {
			S: aws.String(e.Name),
		},
		"value": {
			S: aws.String(e.Value),
		},
	}
	if e.Error != nil {
		er := *e.Error
		errMessage := er.Error()
		logObject = map[string]*ddb.AttributeValue{
			"name": {
				S: aws.String(e.Name),
			},
			"value": {
				S: aws.String(e.Value),
			},
			"error": {
				S: aws.String(errMessage),
			},
		}
	}

	logArray := []*ddb.AttributeValue{
		{
			M: logObject,
		},
	}
	vals := map[string]*ddb.AttributeValue{
		":logs": {
			L: logArray,
		},
	}
	_, err := client.UpdateItem(&ddb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(`SET logs = list_append(logs, :logs)`),
		ExpressionAttributeValues: vals,
	})
	if err != nil {
		errMessage := err.Error()
		awsNotFound := "The provided expression refers to an attribute that does not exist in the item"
		if strings.Contains(errMessage, awsNotFound) {
			_, err := client.UpdateItem(&ddb.UpdateItemInput{
				TableName:                 aws.String(tableName),
				Key:                       key,
				UpdateExpression:          aws.String(`SET logs = :logs`),
				ExpressionAttributeValues: vals,
			})
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}

	}
}
