package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
)

type PlanUsageInfo struct {
	BusinessID        string `json:"business_id"`
	PlanID            string `json:"plan_id"`
	PlanName          string `json:"plan_name"`
	NotificationLimit int    `json:"notification_limit"`
	NotificationCount int    `json:"notification_count"`
	NotificationLeft  int    `json:"notification_left"`
	PeriodStart       string `json:"period_start"`
	PeriodEnd         string `json:"period_end"`
	PeriodDays        int    `json:"period_days"`
}

func GetPlanUsageService(apiKey string) (*PlanUsageInfo, error) {
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

	// Obtener información del plan
	plan, err := planRepo.GetByID(ctx, business.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan not found")
	}

	// Obtener uso actual
	businessID := business.PK[9:] // Remover "BUSINESS#"
	usage, err := usageRepo.GetCurrentUsage(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("usage not found")
	}

	notificationLeft := plan.NotificationLimit - usage.NotificationCount
	if notificationLeft < 0 {
		notificationLeft = 0
	}

	info := &PlanUsageInfo{
		BusinessID:        businessID,
		PlanID:            plan.PK[5:], // Remover "PLAN#"
		PlanName:          plan.Name,
		NotificationLimit: plan.NotificationLimit,
		NotificationCount: usage.NotificationCount,
		NotificationLeft:  notificationLeft,
		PeriodStart:       usage.PeriodStart,
		PeriodEnd:         usage.PeriodEnd,
		PeriodDays:        plan.PeriodDays,
	}

	return info, nil
}
