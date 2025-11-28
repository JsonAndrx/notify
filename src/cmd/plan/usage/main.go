package main

import (
	"notify-backend/common/response"
	"notify-backend/internal/services"
	"notify-backend/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func PlanUsageHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Obtener API Key del header
	apiKey := utils.ExtractAPIKeyFromRequest(request)

	if apiKey == "" {
		return response.ErrorResponse(401, "API Key is required"), nil
	}

	usage, err := services.GetPlanUsageService(apiKey)
	if err != nil {
		return response.ErrorResponse(401, err.Error()), nil
	}

	return response.SuccessResponse(200, usage), nil
}

func main() {
	lambda.Start(PlanUsageHandler)
}
