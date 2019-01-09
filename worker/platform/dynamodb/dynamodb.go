package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jurekbarth/pup/worker"
	"github.com/jurekbarth/pup/worker/internal/project"
)

// Rule for lambda edge server
type Rule struct {
	Groups []string `json:"group-permissions"`
}

// RuleEntry struct
type RuleEntry struct {
	BaseURI string            `json:"baseUri"`
	Rules   []map[string]Rule `json:"r"` // we have to use something else than rules because it's reserved
}

// GetItemByURI is a wrapper returning an item
func GetItemByURI(w *worker.Worker, uri string) (*RuleEntry, error) {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return nil, err
	}
	client := ddb.New(session)
	key := map[string]*ddb.AttributeValue{
		"baseUri": {
			S: aws.String(uri),
		},
	}
	table := w.Config.AWSDynamoDBRulesTable
	res, err := client.GetItem(&ddb.GetItemInput{
		TableName:      &table,
		Key:            key,
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	r := RuleEntry{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// UpsertItemByURI updates or create an item for authtable
func UpsertItemByURI(w *worker.Worker, uri string, r project.Project) error {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return err
	}
	client := ddb.New(session)
	table := w.Config.AWSDynamoDBRulesTable
	key := map[string]*ddb.AttributeValue{
		"baseUri": {
			S: aws.String(uri),
		},
	}
	av, err := dynamodbattribute.Marshal(r.Rules)
	if err != nil {
		return err
	}
	vals := map[string]*ddb.AttributeValue{
		":r": av,
	}
	_, err = client.UpdateItem(&ddb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       key,
		UpdateExpression:          aws.String(`SET r = :r`),
		ExpressionAttributeValues: vals,
	})
	return err
}

// GetAllRules returns all rules
func GetAllRules(w *worker.Worker) (*[]RuleEntry, error) {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return nil, err
	}
	client := ddb.New(session)
	table := w.Config.AWSDynamoDBRulesTable
	res, err := client.Scan(&ddb.ScanInput{
		TableName: &table,
	})
	rules := []RuleEntry{}
	for _, ar := range res.Items {
		r := RuleEntry{}
		err := dynamodbattribute.UnmarshalMap(ar, &r)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return &rules, nil
}
