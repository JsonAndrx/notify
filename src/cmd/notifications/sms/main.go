package main

import (
	"encoding/json"

	"notify-backend/common/response"
	"notify-backend/internal/services"
	"notify-backend/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type SendSMSRequest struct {
	To      string `json:"to" validate:"required"`
	Message string `json:"message" validate:"required,max=1600"`
}

func SendSMSHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Obtener API Key del header
	apiKey := utils.ExtractAPIKeyFromRequest(request)

	if apiKey == "" {
		return response.ErrorResponse(401, "API Key is required"), nil
	}

	var req SendSMSRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.ErrorResponse(400, "Invalid request body"), nil
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return response.ErrorResponse(400, "Invalid request body: "+err.Error()), nil
	}

	serviceReq := services.SendSMSRequest{
		To:      req.To,
		Message: req.Message,
	}

	result, err := services.SendSMSService(apiKey, serviceReq)
	if err != nil {
		statusCode := 500
		errMsg := err.Error()

		if errMsg == "invalid API key" {
			statusCode = 401
		} else if errMsg == "notification limit reached" {
			statusCode = 429
		} else if errMsg == "message cannot be empty" || errMsg == "message too long. Maximum 1600 characters allowed" {
			statusCode = 400
		} else if len(errMsg) > 16 && errMsg[:16] == "failed to send SMS" {
			statusCode = 502
		}

		return response.ErrorResponse(statusCode, errMsg), nil
	}

	return response.SuccessResponse(200, result), nil
}

func main() {
	lambda.Start(SendSMSHandler)
}
