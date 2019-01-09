package cloudfront

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/jurekbarth/pup/worker"
)

// ClearCache ...
func ClearCache(w *worker.Worker, pathPattern string) (*cloudfront.CreateInvalidationOutput, error) {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return nil, err
	}
	cfs := cloudfront.New(session)
	dateStr := strconv.FormatInt(time.Now().Unix(), 10)
	result, err := cfs.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(w.Config.AWSCloudfrontID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(dateStr),
			Paths: &cloudfront.Paths{
				Items:    aws.StringSlice([]string{pathPattern}),
				Quantity: aws.Int64(1),
			},
		},
	})
	return result, err
}

// GetConfig ...
func GetConfig(w *worker.Worker) (*cloudfront.GetDistributionConfigOutput, error) {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return nil, err
	}

	cfs := cloudfront.New(session)
	input := &cloudfront.GetDistributionConfigInput{
		Id: aws.String(w.Config.AWSCloudfrontID),
	}
	return cfs.GetDistributionConfig(input)

}

// UpdateConfig ...
func UpdateConfig(w *worker.Worker, config *cloudfront.GetDistributionConfigOutput, lambdaVersion *string) error {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return err
	}
	cfs := cloudfront.New(session)
	etag := config.ETag
	distributionConfig := config.DistributionConfig
	arn := w.Config.AWSLambdaARN + ":" + *lambdaVersion
	distributionConfig.DefaultCacheBehavior.LambdaFunctionAssociations.Items[0].LambdaFunctionARN = aws.String(arn)
	input := &cloudfront.UpdateDistributionInput{
		Id:                 aws.String(w.Config.AWSCloudfrontID),
		DistributionConfig: distributionConfig,
		IfMatch:            etag,
	}
	_, err = cfs.UpdateDistribution(input)
	return err
}
