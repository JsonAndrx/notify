package repository

import (
	"context"
	"fmt"
	"time"

	"notify-backend/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UsageRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewUsageRepository(client *dynamodb.Client, tableName string) *UsageRepository {
	return &UsageRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (r *UsageRepository) GetCurrentUsage(ctx context.Context, businessID string) (*models.Usage, error) {
	out, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "BUSINESS#" + businessID},
			":sk": &types.AttributeValueMemberS{Value: "USAGE#"},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, fmt.Errorf("usage not found")
	}

	item := out.Items[0]

	notificationCount := 0
	if nc, ok := item["notificationCount"]; ok {
		fmt.Sscanf(nc.(*types.AttributeValueMemberN).Value, "%d", &notificationCount)
	}

	usage := &models.Usage{
		PK:                item["PK"].(*types.AttributeValueMemberS).Value,
		SK:                item["SK"].(*types.AttributeValueMemberS).Value,
		BusinessID:        item["businessId"].(*types.AttributeValueMemberS).Value,
		PlanID:            item["planId"].(*types.AttributeValueMemberS).Value,
		NotificationCount: notificationCount,
		PeriodStart:       item["periodStart"].(*types.AttributeValueMemberS).Value,
		PeriodEnd:         item["periodEnd"].(*types.AttributeValueMemberS).Value,
		CreatedAt:         item["createdAt"].(*types.AttributeValueMemberS).Value,
	}

	if updatedAt, ok := item["updatedAt"]; ok {
		usage.UpdatedAt = updatedAt.(*types.AttributeValueMemberS).Value
	}

	return usage, nil
}

func (r *UsageRepository) Create(ctx context.Context, usage *models.Usage) error {
	item := map[string]types.AttributeValue{
		"PK":                &types.AttributeValueMemberS{Value: usage.PK},
		"SK":                &types.AttributeValueMemberS{Value: usage.SK},
		"businessId":        &types.AttributeValueMemberS{Value: usage.BusinessID},
		"planId":            &types.AttributeValueMemberS{Value: usage.PlanID},
		"notificationCount": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", usage.NotificationCount)},
		"periodStart":       &types.AttributeValueMemberS{Value: usage.PeriodStart},
		"periodEnd":         &types.AttributeValueMemberS{Value: usage.PeriodEnd},
		"createdAt":         &types.AttributeValueMemberS{Value: usage.CreatedAt},
	}

	if usage.UpdatedAt != "" {
		item["updatedAt"] = &types.AttributeValueMemberS{Value: usage.UpdatedAt}
	}

	_, err := r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	})

	return err
}

func (r *UsageRepository) IncrementUsage(ctx context.Context, businessID, usageSK string) error {
	_, err := r.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "BUSINESS#" + businessID},
			"SK": &types.AttributeValueMemberS{Value: usageSK},
		},
		UpdateExpression: aws.String("ADD notificationCount :inc SET updatedAt = :updatedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc":       &types.AttributeValueMemberN{Value: "1"},
			":updatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	})

	return err
}

func (r *UsageRepository) CheckAndCreateNewPeriod(ctx context.Context, businessID, planID string, periodDays int) (*models.Usage, error) {
	// Intentar obtener el uso actual
	currentUsage, err := r.GetCurrentUsage(ctx, businessID)
	if err == nil {
		// Verificar si el período ha expirado
		periodEnd, _ := time.Parse(time.RFC3339, currentUsage.PeriodEnd)
		if time.Now().Before(periodEnd) {
			return currentUsage, nil
		}
	}

	// Crear nuevo período
	now := time.Now()
	periodEnd := now.AddDate(0, 0, periodDays)

	newUsage := &models.Usage{
		PK:                "BUSINESS#" + businessID,
		SK:                "USAGE#" + now.Format("2006-01-02"),
		BusinessID:        businessID,
		PlanID:            planID,
		NotificationCount: 0,
		PeriodStart:       now.Format(time.RFC3339),
		PeriodEnd:         periodEnd.Format(time.RFC3339),
		CreatedAt:         now.Format(time.RFC3339),
	}

	if err := r.Create(ctx, newUsage); err != nil {
		return nil, err
	}

	return newUsage, nil
}
