package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"notify-backend/internal/utils"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SendSMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type SendSMSResponse struct {
	Success           bool   `json:"success"`
	NotificationID    string `json:"notification_id"`
	NotificationCount int    `json:"notification_count"`
	NotificationLeft  int    `json:"notification_left"`
}

func SendSMSService(apiKey string, req SendSMSRequest) (*SendSMSResponse, error) {
	client, _ := db.NewDynamoClient()
	businessRepo := repository.NewBusinessRepository(client, "NotificationService")
	planRepo := repository.NewPlanRepository(client, "NotificationService")
	usageRepo := repository.NewUsageRepository(client, "NotificationService")
	ctx := context.TODO()

	// Validar formato del número de teléfono
	formattedPhone := utils.FormatPhoneNumber(req.To)
	if !utils.ValidatePhoneNumber(formattedPhone) {
		return nil, fmt.Errorf("invalid phone number format")
	}
	req.To = formattedPhone

	// Validar longitud del mensaje
	if len(req.Message) == 0 {
		return nil, fmt.Errorf("message cannot be empty")
	}

	if len(req.Message) > 1600 {
		return nil, fmt.Errorf("message too long")
	}

	// Buscar negocio por API Key
	business, err := businessRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("authentication failed")
	}

	businessID := business.PK[9:] // Remover "BUSINESS#"

	// Obtener plan
	plan, err := planRepo.GetByID(ctx, business.PlanID)
	if err != nil {
		return nil, fmt.Errorf("service unavailable")
	}

	// Verificar o crear período de uso
	usage, err := usageRepo.CheckAndCreateNewPeriod(ctx, businessID, business.PlanID, plan.PeriodDays)
	if err != nil {
		return nil, fmt.Errorf("service unavailable")
	}

	// Verificar límite de notificaciones
	if usage.NotificationCount >= plan.NotificationLimit {
		return nil, fmt.Errorf("notification limit reached")
	}

	// Obtener cliente de Twilio
	twilioClient := GetTwilioClient()
	twilioPhoneNumber := GetTwilioPhoneNumber()

	// Verificar que Twilio esté configurado
	if twilioClient == nil || twilioPhoneNumber == "" {
		return nil, fmt.Errorf("service temporarily unavailable")
	}

	// Enviar SMS a través de Twilio
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(req.To)
	params.SetFrom(twilioPhoneNumber)
	params.SetBody(req.Message)

	message, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		// Log interno del error real para debugging
		fmt.Printf("Twilio SMS error: %v\n", err)
		return nil, fmt.Errorf("failed to send notification")
	}

	notificationID := *message.Sid
	fmt.Printf("SMS sent - MessageSID: %s, To: %s\n", notificationID, req.To)

	// Incrementar contador de uso
	err = usageRepo.IncrementUsage(ctx, businessID, usage.SK)
	if err != nil {
		// Log interno para debugging
		fmt.Printf("Failed to increment usage: %v\n", err)
		// No fallar la request si el mensaje ya fue enviado
	}

	// Calcular notificaciones restantes
	newCount := usage.NotificationCount + 1
	notificationLeft := plan.NotificationLimit - newCount
	if notificationLeft < 0 {
		notificationLeft = 0
	}

	return &SendSMSResponse{
		Success:           true,
		NotificationID:    notificationID,
		NotificationCount: newCount,
		NotificationLeft:  notificationLeft,
	}, nil
}
