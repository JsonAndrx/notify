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

type SendEmailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
	HTML    bool   `json:"html"`
}

func SendEmailHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Obtener API Key del header
	apiKey := utils.ExtractAPIKeyFromRequest(request)

	if apiKey == "" {
		return response.ErrorResponse(401, "API Key is required"), nil
	}

	var req SendEmailRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.ErrorResponse(400, "Invalid request body"), nil
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return response.ErrorResponse(400, "Invalid request body: "+err.Error()), nil
	}

	serviceReq := services.SendEmailRequest{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
		HTML:    req.HTML,
	}

	result, err := services.SendEmailService(apiKey, serviceReq)
	if err != nil {
		statusCode := 500
		errMsg := err.Error()

		if errMsg == "invalid API key" {
			statusCode = 401
		} else if errMsg == "notification limit reached" {
			statusCode = 429
		}

		return response.ErrorResponse(statusCode, errMsg), nil
	}

	return response.SuccessResponse(200, result), nil
}

func main() {
	lambda.Start(SendEmailHandler)
}
