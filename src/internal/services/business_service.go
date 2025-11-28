package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/models"
	"notify-backend/internal/repository"
	"notify-backend/internal/utils"
	"time"

	"github.com/google/uuid"
)

type BusinessRegisterResult struct {
	BusinessID string
	APIKey     string
}

func BusinessRegisterService(name string, email string, phone string, planID string) (*BusinessRegisterResult, error) {
	client, _ := db.NewDynamoClient()
	repo := repository.NewBusinessRepository(client, "NotificationService")
	planRepo := repository.NewPlanRepository(client, "NotificationService")
	usageRepo := repository.NewUsageRepository(client, "NotificationService")
	ctx := context.TODO()

	// Si no se especifica plan, usar FREE por defecto
	if planID == "" {
		planID = "FREE"
	}

	// Verificar que el plan existe
	plan, err := planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan")
	}

	// Verificar email existente
	emailExists, err := repo.EmailExists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("service unavailable")
	}
	if emailExists {
		return nil, fmt.Errorf("email already registered")
	}

	// Verificar phone existente
	phoneExists, err := repo.PhoneExists(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("service unavailable")
	}
	if phoneExists {
		return nil, fmt.Errorf("phone already registered")
	}

	// Generar API Key
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("service unavailable")
	}

	// Crear nuevo negocio
	id := uuid.New().String()
	item := &models.Business{
		PK:        "BUSINESS#" + id,
		SK:        "METADATA",
		Name:      name,
		Email:     email,
		Phone:     phone,
		PlanID:    planID,
		APIKey:    apiKey,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := repo.Create(ctx, item); err != nil {
		return nil, err
	}

	// Crear período de uso inicial
	now := time.Now()
	periodEnd := now.AddDate(0, 0, plan.PeriodDays)

	usage := &models.Usage{
		PK:                "BUSINESS#" + id,
		SK:                "USAGE#" + now.Format("2006-01-02"),
		BusinessID:        id,
		PlanID:            planID,
		NotificationCount: 0,
		PeriodStart:       now.Format(time.RFC3339),
		PeriodEnd:         periodEnd.Format(time.RFC3339),
		CreatedAt:         now.Format(time.RFC3339),
	}

	if err := usageRepo.Create(ctx, usage); err != nil {
		return nil, fmt.Errorf("service unavailable")
	}

	return &BusinessRegisterResult{
		BusinessID: id,
		APIKey:     apiKey,
	}, nil
}
