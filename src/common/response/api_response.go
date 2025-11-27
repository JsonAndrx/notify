package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type StructResp struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

func SuccessResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	resp := StructResp{
		Status: true,
		Data:   data,
		Error:  "",
	}
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func ErrorResponse(statusCode int, errorRes interface{}) events.APIGatewayProxyResponse {
	resp := StructResp{
		Status: false,
		Data:   nil,
		Error:  errorRes,
	}
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
