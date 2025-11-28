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

type SendNotificationRequest struct {
	To      string `json:"to" validate:"required"`
	Message string `json:"message" validate:"required"`
	Type    string `json:"type" validate:"required,oneof=whatsapp sms email"`
}

func SendNotificationHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Obtener API Key del header
	apiKey := utils.ExtractAPIKeyFromRequest(request)

	if apiKey == "" {
		return response.ErrorResponse(401, "API Key is required"), nil
	}

	var req SendNotificationRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.ErrorResponse(400, "Invalid request body"), nil
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return response.ErrorResponse(400, "Invalid request body: "+err.Error()), nil
	}

	serviceReq := services.SendNotificationRequest{
		To:      req.To,
		Message: req.Message,
		Type:    req.Type,
	}

	result, err := services.SendNotificationService(apiKey, serviceReq)
	if err != nil {
		statusCode := 500
		if err.Error() == "invalid API key" {
			statusCode = 401
		} else if err.Error() == "notification limit reached. Please upgrade your plan" {
			statusCode = 429
		}
		return response.ErrorResponse(statusCode, err.Error()), nil
	}

	return response.SuccessResponse(200, result), nil
}

func main() {
	lambda.Start(SendNotificationHandler)
}
