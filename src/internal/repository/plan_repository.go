package repository

import (
	"context"
	"fmt"

	"notify-backend/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PlanRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewPlanRepository(client *dynamodb.Client, tableName string) *PlanRepository {
	return &PlanRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (r *PlanRepository) GetByID(ctx context.Context, planID string) (*models.Plan, error) {
	out, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "PLAN#" + planID},
			"SK": &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, fmt.Errorf("plan not found")
	}

	var notificationLimit, periodDays int
	var price float64

	fmt.Sscanf(out.Item["notificationLimit"].(*types.AttributeValueMemberN).Value, "%d", &notificationLimit)
	fmt.Sscanf(out.Item["periodDays"].(*types.AttributeValueMemberN).Value, "%d", &periodDays)

	plan := &models.Plan{
		PK:                out.Item["PK"].(*types.AttributeValueMemberS).Value,
		SK:                out.Item["SK"].(*types.AttributeValueMemberS).Value,
		Name:              out.Item["name"].(*types.AttributeValueMemberS).Value,
		NotificationLimit: notificationLimit,
		PeriodDays:        periodDays,
		Description:       out.Item["description"].(*types.AttributeValueMemberS).Value,
		Active:            out.Item["active"].(*types.AttributeValueMemberBOOL).Value,
		CreatedAt:         out.Item["createdAt"].(*types.AttributeValueMemberS).Value,
	}

	if priceItem, ok := out.Item["price"]; ok {
		fmt.Sscanf(priceItem.(*types.AttributeValueMemberN).Value, "%f", &price)
		plan.Price = price
	}

	return plan, nil
}

func (r *PlanRepository) Create(ctx context.Context, plan *models.Plan) error {
	item := map[string]types.AttributeValue{
		"PK":                &types.AttributeValueMemberS{Value: plan.PK},
		"SK":                &types.AttributeValueMemberS{Value: plan.SK},
		"name":              &types.AttributeValueMemberS{Value: plan.Name},
		"notificationLimit": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", plan.NotificationLimit)},
		"periodDays":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", plan.PeriodDays)},
		"price":             &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", plan.Price)},
		"description":       &types.AttributeValueMemberS{Value: plan.Description},
		"active":            &types.AttributeValueMemberBOOL{Value: plan.Active},
		"createdAt":         &types.AttributeValueMemberS{Value: plan.CreatedAt},
	}

	_, err := r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	})

	return err
}
