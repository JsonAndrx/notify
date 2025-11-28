package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"notify-backend/internal/utils"
	"time"
)

func RegenerateAPIKeyService(email, phone, currentAPIKey string) (string, error) {
	client, _ := db.NewDynamoClient()
	repo := repository.NewBusinessRepository(client, "NotificationService")
	ctx := context.TODO()

	// Buscar negocio por API Key actual
	business, err := repo.GetByAPIKey(ctx, currentAPIKey)
	if err != nil {
		return "", fmt.Errorf("business not found")
	}

	// Verificar que el email y phone coincidan
	if business.Email != email || business.Phone != phone {
		return "", fmt.Errorf("invalid credentials")
	}

	// Generar nueva API Key
	newAPIKey, err := utils.GenerateAPIKey()
	if err != nil {
		return "", fmt.Errorf("service unavailable")
	}

	// Actualizar en base de datos
	err = repo.UpdateAPIKey(ctx, business.PK, business.APIKey, newAPIKey, time.Now().Format(time.RFC3339))
	if err != nil {
		return "", fmt.Errorf("service unavailable")
	}

	return newAPIKey, nil
}

func GetBusinessInfoService(apiKey string) (*BusinessInfo, error) {
	client, _ := db.NewDynamoClient()
	repo := repository.NewBusinessRepository(client, "NotificationService")
	ctx := context.TODO()

	// Buscar negocio por API Key
	business, err := repo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("authentication failed")
	}

	info := &BusinessInfo{
		IDBusiness: business.PK[9:], // Remover "BUSINESS#"
		Name:       business.Name,
		Email:      business.Email,
		Phone:      business.Phone,
		PlanID:     business.PlanID,
		CreatedAt:  business.CreatedAt,
	}

	return info, nil
}

type BusinessInfo struct {
	IDBusiness string `json:"id_business"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	PlanID     string `json:"plan_id"`
	CreatedAt  string `json:"created_at"`
}
