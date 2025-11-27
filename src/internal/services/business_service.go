package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/models"
	"notify-backend/internal/repository"
	"time"

	"github.com/google/uuid"
)

func BusinessRegisterService(name string, email string, phone string, planID string) (string, error) {
    client, _ := db.NewDynamoClient()
    repo := repository.NewBusinessRepository(client, "NotificationService")
    ctx := context.TODO()

    // Verificar email existente
    emailExists, err := repo.EmailExists(ctx, email)
    if err != nil {
        return "", err
    }
    if emailExists {
        return "", fmt.Errorf("email already registered")
    }

    // Verificar phone existente
    phoneExists, err := repo.PhoneExists(ctx, phone)
    if err != nil {
        return "", err
    }
    if phoneExists {
        return "", fmt.Errorf("phone already registered")
    }

    // Crear nuevo
    id := uuid.New().String()
    item := &models.Business{
        PK:        "BUSINESS#" + id,
        SK:        "METADATA",
        Name:      name,
        Email:     email,
        Phone:     phone,
        PlanID:    planID,
        CreatedAt: time.Now().Format(time.RFC3339),
    }

    if err := repo.Create(ctx, item); err != nil {
        return "", err
    }

    return id, nil
}