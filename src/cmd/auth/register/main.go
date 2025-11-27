package main

import (
	"encoding/json"

	"notify-backend/common/response"
	"notify-backend/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	PlanID     string `json:"plan_id" validate:"required"`
}

type RegisterResponse struct {
	IDBusiness string `json:"id_business"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	PlanID     string `json:"plan_id"`
}

func RegisterBusinessHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var req RegisterRequest
    if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
        return response.ErrorResponse(400, "Invalid request body"), nil
    }

	validate := validator.New()
    if err := validate.Struct(req); err != nil {
        return response.ErrorResponse(400, "Invalid request body"), nil
    }

    idBusiness, err := services.BusinessRegisterService(
        req.Name,
        req.Email,
        req.Phone,
        req.PlanID,
    )

    if err != nil {
        return response.ErrorResponse(500, err.Error()), nil
    }

    resp := RegisterResponse{
        IDBusiness: idBusiness,
        Name:       req.Name,
        Email:      req.Email,
        Phone:      req.Phone,
        PlanID:     req.PlanID,
    }

    return response.SuccessResponse(200, resp), nil
}

func main() {
	lambda.Start(RegisterBusinessHandler)
}