package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"time"

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

	// Validar longitud del mensaje
	if len(req.Message) == 0 {
		return nil, fmt.Errorf("message cannot be empty")
	}

	// SMS estándar: 160 caracteres para GSM-7, 70 para Unicode
	// Permitimos hasta 1600 caracteres (10 segmentos concatenados)
	if len(req.Message) > 1600 {
		return nil, fmt.Errorf("message too long. Maximum 1600 characters allowed")
	}

	// Buscar negocio por API Key
	business, err := businessRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}

	businessID := business.PK[9:] // Remover "BUSINESS#"

	// Obtener plan
	plan, err := planRepo.GetByID(ctx, business.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan not found")
	}

	// Verificar o crear período de uso
	usage, err := usageRepo.CheckAndCreateNewPeriod(ctx, businessID, business.PlanID, plan.PeriodDays)
	if err != nil {
		return nil, fmt.Errorf("failed to check usage: %v", err)
	}

	// Verificar límite de notificaciones
	if usage.NotificationCount >= plan.NotificationLimit {
		return nil, fmt.Errorf("notification limit reached")
	}

	var notificationID string

	// Obtener cliente de Twilio
	twilioClient := GetTwilioClient()
	twilioPhoneNumber := GetTwilioPhoneNumber()

	if twilioClient != nil && twilioPhoneNumber != "" {
		// Integración real con Twilio SMS
		params := &twilioApi.CreateMessageParams{}
		params.SetTo(req.To)
		params.SetFrom(twilioPhoneNumber)
		params.SetBody(req.Message)

		message, err := twilioClient.Api.CreateMessage(params)
		if err != nil {
			return nil, fmt.Errorf("failed to send SMS: %v", err)
		}

		notificationID = *message.Sid
		fmt.Printf("SMS sent - MessageSID: %s, To: %s\n", notificationID, req.To)
	} else {
		// Modo simulación (sin credenciales de Twilio)
		notificationID = fmt.Sprintf("SMS_SIM_%d", time.Now().UnixNano())
		fmt.Printf("SMS SIMULATION - To: %s, Message: %s\n", req.To, req.Message)
		fmt.Println("💡 Tip: Configure TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN y TWILIO_PHONE_NUMBER para envíos reales")
	}

	// Incrementar contador de uso
	err = usageRepo.IncrementUsage(ctx, businessID, usage.SK)
	if err != nil {
		return nil, fmt.Errorf("failed to increment usage: %v", err)
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
