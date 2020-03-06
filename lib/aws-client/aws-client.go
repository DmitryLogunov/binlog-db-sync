package awsClient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"binlog-db-sync/lib/types"
	"binlog-db-sync/lib/logger"
	"encoding/json"
	"strconv"
)

type InvokeResponse struct {
	StatusCode int
	Body       interface{}
}

func Invoke(functionName string, params *types.InvokeParams) (map[string]interface{}) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	payload, err := json.Marshal(params)
	if err != nil {
		logger.Error("invoke response: error marshalling params for invoke aws-lambda-function", map[string]interface{}{"lambda-function": functionName, "params": params})
		return map[string]interface{}{"indo": "Error: invoke response: error marshalling params for invoke aws-lambda-function", "lambda-function": functionName, "params": params}
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(functionName), Payload: payload})
	if err != nil {
		logger.Error("invoke response: error calling  aws-lambda-function", map[string]interface{}{"lambda-function": functionName, "params": params})
		return map[string]interface{}{"indo": "Error: invoke response: error calling  aws-lambda-function", "lambda-function": functionName, "params": params, "invokeResult": result}
	}

	resp := InvokeResponse{}
	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		logger.Error("invoke response: error unmarshalling aws-lambda-function", map[string]interface{}{"lambda-function": functionName, "params": params})
		return map[string]interface{}{"indo": "Error: invoke response: error calling  aws-lambda-function", "lambda-function": functionName, "params": params, "response": resp}
	}

	if resp.StatusCode != 200 {
		logger.Error("invoke response:  error of invoking aws-lambda-function response", map[string]interface{}{"lambda-function": functionName, "params": params, "status-code": strconv.Itoa(resp.StatusCode)})
		return map[string]interface{}{"indo": "Error: invoke response: error calling  aws-lambda-function", "lambda-function": functionName, "params": params, "response": resp}
	}

	return map[string]interface{}{"info": "Success", "lambda-function": functionName, "params": params, "response": resp}
}
