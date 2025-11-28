package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"time"
)

type SendNotificationRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Type    string `json:"type"` // whatsapp, sms, email
}

type SendNotificationResponse struct {
	Success           bool   `json:"success"`
	NotificationID    string `json:"notification_id"`
	NotificationCount int    `json:"notification_count"`
	NotificationLeft  int    `json:"notification_left"`
}

func SendNotificationService(apiKey string, req SendNotificationRequest) (*SendNotificationResponse, error) {
	client, _ := db.NewDynamoClient()
	businessRepo := repository.NewBusinessRepository(client, "NotificationService")
	planRepo := repository.NewPlanRepository(client, "NotificationService")
	usageRepo := repository.NewUsageRepository(client, "NotificationService")
	ctx := context.TODO()

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
		return nil, fmt.Errorf("notification limit reached. Please upgrade your plan")
	}

	// TODO: Aquí iría la lógica real para enviar la notificación
	// Por ahora solo simularemos el envío
	notificationID := fmt.Sprintf("NOTIF_%d", time.Now().UnixNano())

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

	return &SendNotificationResponse{
		Success:           true,
		NotificationID:    notificationID,
		NotificationCount: newCount,
		NotificationLeft:  notificationLeft,
	}, nil
}
