package main

import (
	"encoding/json"

	"notify-backend/common/response"
	"notify-backend/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type RegenerateKeyRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Phone         string `json:"phone" validate:"required"`
	CurrentAPIKey string `json:"current_api_key" validate:"required"`
}

type RegenerateKeyResponse struct {
	APIKey string `json:"api_key"`
}

func RegenerateKeyHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req RegenerateKeyRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.ErrorResponse(400, "Invalid request body"), nil
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return response.ErrorResponse(400, "Invalid request body: "+err.Error()), nil
	}

	newAPIKey, err := services.RegenerateAPIKeyService(req.Email, req.Phone, req.CurrentAPIKey)
	if err != nil {
		statusCode := 500
		if err.Error() == "invalid credentials" || err.Error() == "business not found" {
			statusCode = 401
		}
		return response.ErrorResponse(statusCode, err.Error()), nil
	}

	resp := RegenerateKeyResponse{
		APIKey: newAPIKey,
	}

	return response.SuccessResponse(200, resp), nil
}

func main() {
	lambda.Start(RegenerateKeyHandler)
}
