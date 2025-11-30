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
	To         string            `json:"to" validate:"required"`
	TemplateID string            `json:"template_id" validate:"required"`
	Parameters map[string]string `json:"parameters"`
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
		return response.ErrorResponse(400, "Invalid request body"), nil
	}

	serviceReq := services.SendSMSRequest{
		To:         req.To,
		TemplateID: req.TemplateID,
		Parameters: req.Parameters,
	}

	result, err := services.SendSMSService(apiKey, serviceReq)
	if err != nil {
		statusCode := 500
		errMsg := err.Error()

		if errMsg == "authentication failed" {
			statusCode = 401
		} else if errMsg == "notification limit reached" {
			statusCode = 429
		} else if errMsg == "invalid phone number format" || errMsg == "invalid template" || errMsg == "invalid template type" || errMsg == "missing required parameters" || errMsg == "invalid verification code format" || errMsg == "message too long" {
			statusCode = 400
		} else if errMsg == "template not available" {
			statusCode = 404
		}

		return response.ErrorResponse(statusCode, errMsg), nil
	}

	return response.SuccessResponse(200, result), nil
}

func main() {
	lambda.Start(SendSMSHandler)
}
