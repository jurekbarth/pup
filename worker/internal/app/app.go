package app

import (
	"fmt"
	"github.com/jurekbarth/pup/worker"
	"github.com/jurekbarth/pup/worker/event"
	"github.com/jurekbarth/pup/worker/internal/project"
	"github.com/jurekbarth/pup/worker/internal/unzip"
	"github.com/jurekbarth/pup/worker/internal/usermanagement"
	"github.com/jurekbarth/pup/worker/platform/cloudfront"
	"github.com/jurekbarth/pup/worker/platform/dynamodb"
	"github.com/jurekbarth/pup/worker/platform/lambda"
	"github.com/jurekbarth/pup/worker/platform/s3"
	"github.com/jurekbarth/pup/worker/platform/sqs"
	"github.com/jurekbarth/pup/worker/reporter"
	"github.com/pkg/errors"
)

func compareRules(dbRuleSet []map[string]dynamodb.Rule, configRuleSet []map[string]project.Rule) bool {
	for idx, dbRules := range dbRuleSet {
		for key, dbRule := range dbRules {
			configRule := configRuleSet[idx][key]
			if len(configRule.Groups) != len(dbRule.Groups) {
				return true
			}
			for i, group := range dbRule.Groups {
				if group != configRule.Groups[i] {
					return true
				}
			}
		}
	}
	return false
}

// Run the steps
func Run() error {
	_, w, err := start()
	sqsMessage, err := sqs.GetSqsMessage(w)
	if err != nil {
		return err
	}
	key, err := sqs.GetS3Key(sqsMessage)
	if err != nil {
		return err
	}
	// key := "1546928870-vi-test-project-latest.zip"
	l := len([]rune(key))
	id := key[0 : l-4]
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "s3download.start",
		Value: "download zip file",
		Error: nil,
	})
	err = s3.Download(w, key)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "s3download.error",
			Value: "err",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "s3download.done",
		Value: "download done",
		Error: nil,
	})
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "unzip.start",
		Value: "unzip file",
		Error: nil,
	})
	err = unzip.Unzip(w, key)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "unzip.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "unzip.pending",
		Value: "read config",
		Error: nil,
	})
	pConfig, dConfig, err := project.Read(w)
	fmt.Println("############ pConfig ############")
	fmt.Println(pConfig)
	fmt.Println("############ dConfig ############")
	fmt.Println(dConfig)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "unzip.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "unzip.done",
		Value: "successful unzipped and read",
		Error: nil,
	})
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "usermanagement.start",
		Value: "setting up users",
		Error: nil,
	})
	err = usermanagement.Create(w, pConfig.Users)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "usermanagement.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "usermanagement.done",
		Value: "setting up users",
		Error: nil,
	})
	rulesPath := "/" + pConfig.CustomerCode + "/" + pConfig.Project
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "ddb.start",
		Value: "get dynamo db key",
		Error: nil,
	})
	rulesEntry, err := dynamodb.GetItemByURI(w, rulesPath)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "ddb.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "ddb.done",
		Value: "fetched data",
		Error: nil,
	})

	configRuleSet := pConfig.Rules
	dbRuleSet := rulesEntry.Rules

	lambdaNeedsUpdate := false
	if len(configRuleSet) != len(dbRuleSet) {
		lambdaNeedsUpdate = true
	}
	if !lambdaNeedsUpdate {
		lambdaNeedsUpdate = compareRules(dbRuleSet, configRuleSet)
	}

	if lambdaNeedsUpdate {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "ddb.start",
			Value: "update dynamo db",
			Error: nil,
		})
		err = dynamodb.UpsertItemByURI(w, rulesPath, *pConfig)
		if err != nil {
			w.Events.Emit(event.Event{
				ID:    id,
				Name:  "ddb.error",
				Value: "error",
				Error: &err,
			})
			return err
		}
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "ddb.pending",
			Value: "get all accessrules",
			Error: nil,
		})
		dbRuleEntries, err := dynamodb.GetAllRules(w)
		if err != nil {
			w.Events.Emit(event.Event{
				ID:    id,
				Name:  "ddb.error",
				Value: "error",
				Error: &err,
			})
			return err
		}
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "ddb.done",
			Value: "got all accessrules",
			Error: nil,
		})
		lambdaZipPath := "./" + dConfig.ID + "lambda.zip"
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "lambda.start",
			Value: "generate lambda zip",
			Error: nil,
		})
		err = lambda.GenerateZip(*w.Config, dbRuleEntries, lambdaZipPath)
		if err != nil {
			w.Events.Emit(event.Event{
				ID:    id,
				Name:  "lambda.error",
				Value: "error",
				Error: &err,
			})
			return err
		}
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "lambda.pending",
			Value: "upload lambda zip",
			Error: nil,
		})
		lambdaOutput, err := lambda.UpdateFunction(w, lambdaZipPath)
		if err != nil {
			w.Events.Emit(event.Event{
				ID:    id,
				Name:  "lambda.error",
				Value: "error",
				Error: &err,
			})
			return err
		}
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "lambda.done",
			Value: "uploaded lambda",
			Error: nil,
		})
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "cloudfrontupdate.start",
			Value: "get cloudfront config",
			Error: nil,
		})
		cloudfrontConfig, err := cloudfront.GetConfig(w)
		if err != nil {
			w.Events.Emit(event.Event{
				ID:    id,
				Name:  "cloudfrontupdate",
				Value: "error",
				Error: &err,
			})
			return err
		}
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "cloudfrontupdate.pending",
			Value: "update cloudfront config",
			Error: nil,
		})
		lambdaVersion := lambdaOutput.Version
		err = cloudfront.UpdateConfig(w, cloudfrontConfig, lambdaVersion)
		if err != nil {
			w.Events.Emit(event.Event{
				ID:    id,
				Name:  "cloudfrontupdate.error",
				Value: "error",
				Error: &err,
			})
			return err
		}
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "cloudfrontupdate.done",
			Value: "cloudfrontid: " + w.Config.AWSCloudfrontID,
			Error: nil,
		})
	}
	uploadDestination := rulesPath[1:] + "/" + dConfig.Stage
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "s3.start",
		Value: "clean directory",
		Error: nil,
	})
	err = s3.DeleteDir(w, uploadDestination)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "s3.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "s3.pending",
		Value: "upload started",
		Error: nil,
	})
	fmt.Println("############ uploadDestination ############")
	fmt.Println(uploadDestination)
	err = s3.Upload(w, uploadDestination)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "s3.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "s3.done",
		Value: "upload finished",
		Error: nil,
	})
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "cloudfront.start",
		Value: "clear cache",
		Error: nil,
	})
	cloudfrontInvalidation, err := cloudfront.ClearCache(w, "/"+uploadDestination+"/*")
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "cloudfront.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	invalidationID := ""
	if cloudfrontInvalidation.Invalidation.Id != nil {
		invalidationID = *cloudfrontInvalidation.Invalidation.Id
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "cloudfront.done",
		Value: "cache cleared: " + invalidationID + "###" + w.Config.AWSCloudfrontID,
		Error: nil,
	})
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "sqs.start",
		Value: "delete sqs message",
		Error: nil,
	})
	_, err = sqs.DeleteSqsMessage(w, sqsMessage)
	if err != nil {
		w.Events.Emit(event.Event{
			ID:    id,
			Name:  "sqs.error",
			Value: "error",
			Error: &err,
		})
		return err
	}
	w.Events.Emit(event.Event{
		ID:    id,
		Name:  "sqs.done",
		Value: "success",
		Error: nil,
	})
	// in case return latest error

	return err

}

func start() (*worker.Config, *worker.Worker, error) {
	c, err := worker.ReadConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading config")
	}
	events := make(event.Events)
	w := worker.New(c, events)
	go reporter.DB(events, w)

	return c, w, nil
}
