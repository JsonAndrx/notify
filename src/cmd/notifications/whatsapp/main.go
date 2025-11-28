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

type SendWhatsAppRequest struct {
	To         string            `json:"to" validate:"required"`
	TemplateID string            `json:"template_id" validate:"required"`
	Parameters map[string]string `json:"parameters"`
}

func SendWhatsAppHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Obtener API Key del header
	apiKey := utils.ExtractAPIKeyFromRequest(request)

	if apiKey == "" {
		return response.ErrorResponse(401, "API Key is required"), nil
	}

	var req SendWhatsAppRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.ErrorResponse(400, "Invalid request body"), nil
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return response.ErrorResponse(400, "Invalid request body: "+err.Error()), nil
	}

	serviceReq := services.SendWhatsAppRequest{
		To:         req.To,
		TemplateID: req.TemplateID,
		Parameters: req.Parameters,
	}

	result, err := services.SendWhatsAppService(apiKey, serviceReq)
	if err != nil {
		statusCode := 500
		errMsg := err.Error()

		if errMsg == "invalid API key" {
			statusCode = 401
		} else if errMsg == "notification limit reached" {
			statusCode = 429
		} else if errMsg == "template not found" || errMsg == "template is not active" {
			statusCode = 404
		} else if len(errMsg) > 20 && errMsg[:20] == "invalid template type" {
			statusCode = 400
		} else if len(errMsg) > 25 && errMsg[:25] == "missing required parameters" {
			statusCode = 400
		}

		return response.ErrorResponse(statusCode, errMsg), nil
	}

	return response.SuccessResponse(200, result), nil
}

func main() {
	lambda.Start(SendWhatsAppHandler)
}
