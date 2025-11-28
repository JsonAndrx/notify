package main

import (
	"notify-backend/common/response"
	"notify-backend/internal/services"
	"notify-backend/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func AccountInfoHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Obtener API Key del header
	apiKey := utils.ExtractAPIKeyFromRequest(request)

	if apiKey == "" {
		return response.ErrorResponse(401, "API Key is required"), nil
	}

	info, err := services.GetBusinessInfoService(apiKey)
	if err != nil {
		return response.ErrorResponse(401, "Invalid API Key"), nil
	}

	return response.SuccessResponse(200, info), nil
}

func main() {
	lambda.Start(AccountInfoHandler)
}
