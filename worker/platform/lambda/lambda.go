package lambda

import (
	"fmt"
	"net/http"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	awsLambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jurekbarth/pup/worker"

	"github.com/jurekbarth/pup/worker/internal/zip"
	"github.com/jurekbarth/pup/worker/platform/dynamodb"
	"github.com/jurekbarth/pup/worker/platform/s3"
)

func generateRulesJS(rules *[]dynamodb.RuleEntry) string {
	known := "const known = {"
	for _, r := range *rules {
		str := "'" + r.BaseURI + "': {"
		innerRules := ""
		for _, wrapper := range r.Rules {
			s := ""
			for key, ir := range wrapper {
				st := "'"+key+"':{'groups':["
				groups := ""
				for _, group := range ir.Groups {
					groups = groups + "'"+group+"',"
				}
				en := "],},"
				s = st + groups + en
			}
			innerRules = innerRules + s
		}
		str = str + innerRules + "},"
		known = known + str
	}
	known = known + "}; module.exports = known;"
	return known
}

func generateClientID(clientID string) string {
	return "const clientId = '"+clientID+"'; module.exports = clientId;"
}

func writeStringToFile(filepath string, s string) error {
	bytes := []byte(s)
	err := ioutil.WriteFile(filepath, bytes, 0644)
	return err
}

func getKeysFromAWS(region string, cognitoPoolID string) (string, error) {
	uri := fmt.Sprintf("https://cognito-idp.%v.amazonaws.com/%v/.well-known/jwks.json", region, cognitoPoolID)
	response, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	keys := fmt.Sprintf("%s", string(contents))
	return keys, nil
}

func generateSettings(c worker.Config) (string, error) {
	keys, err := getKeysFromAWS(c.AWSRegion, c.AWSCognitoPoolID)
	if err != nil {
		return "", err
	}
	settings := fmt.Sprintf(`const keys = %s;
const region = "%s";
const cognitoPoolId = "%s";
const endpoint = "%s";
const cfdomain = "%s";
module.exports = {
	keys,
	region,
	cognitoPoolId,
	endpoint,
	cfdomain
}`, keys, c.AWSRegion, c.AWSCognitoPoolID, c.AWSCognitoClientSPASubDomain, c.AWSCloudfrontDomain)
	return settings, nil
}

// GenerateZip generates a zip-file
func GenerateZip(c worker.Config, rules *[]dynamodb.RuleEntry, destination string) error {
	// generates rules.js
	ruleStr := generateRulesJS(rules)
	err := writeStringToFile("./pupauth/rules.js", ruleStr)
	if err != nil {
		return err
	}

	// generate clientId.js
	clientIDStr := generateClientID(c.AWSCognitoClientSPAClientID)
	err = writeStringToFile("./pupauth/clientId.js", clientIDStr)
	if err != nil {
		return err
	}

	// generate settings.js
	settings, err := generateSettings(c)
	if err != nil {
		return err
	}
	err = writeStringToFile("./pupauth/settings.js", settings)
	if err != nil {
		return err
	}

	err = zip.Zip("pupauth", destination)
	return err
}

// UpdateFunction updates a lambda function
func UpdateFunction(w *worker.Worker, filePath string) (*awsLambda.FunctionConfiguration, error) {
	region := "us-east-1"
	session, err := worker.MakeSession(w, &region)
	if err != nil {
		return nil, err
	}
	bucket := w.Config.AWSLambdaBucket
	// remove './' to get a filename out of a path
	filename := filePath[2:]
	err = s3.UploadFileToS3(filePath, session, bucket, "", filename)
	if err != nil {
		return nil, err
	}
	svc := awsLambda.New(session)
	input := &awsLambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(w.Config.AWSLambdaARN),
		Publish:      aws.Bool(true),
		S3Bucket:     aws.String(bucket),
		S3Key:        aws.String(filename),
	}
	return svc.UpdateFunctionCode(input)
}
