package repository

import (
	"context"
	"fmt"

	"notify-backend/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BusinessRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewBusinessRepository(client *dynamodb.Client, tableName string) *BusinessRepository {
	return &BusinessRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (r *BusinessRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	out, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "EMAIL#" + email},
		},
	})
	if err != nil {
		return false, err
	}

	if len(out.Items) == 0 {
		return false, nil
	}

	return true, nil
}

func (r *BusinessRepository) PhoneExists(ctx context.Context, phone string) (bool, error) {
	out, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "PHONE#" + phone},
		},
	})
	if err != nil {
		return false, err
	}

	if len(out.Items) == 0 {
		return false, nil
	}

	return true, nil
}

func (r *BusinessRepository) Create(ctx context.Context, b *models.Business) error {
	emailItem := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "EMAIL#" + b.Email},
		"SK": &types.AttributeValueMemberS{Value: b.PK},
	}

	phoneItem := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "PHONE#" + b.Phone},
		"SK": &types.AttributeValueMemberS{Value: b.PK},
	}

	apiKeyItem := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "APIKEY#" + b.APIKey},
		"SK": &types.AttributeValueMemberS{Value: b.PK},
	}

	metaItem := map[string]types.AttributeValue{
		"PK":        &types.AttributeValueMemberS{Value: b.PK},
		"SK":        &types.AttributeValueMemberS{Value: "METADATA"},
		"name":      &types.AttributeValueMemberS{Value: b.Name},
		"email":     &types.AttributeValueMemberS{Value: b.Email},
		"phone":     &types.AttributeValueMemberS{Value: b.Phone},
		"planId":    &types.AttributeValueMemberS{Value: b.PlanID},
		"apiKey":    &types.AttributeValueMemberS{Value: b.APIKey},
		"createdAt": &types.AttributeValueMemberS{Value: b.CreatedAt},
	}

	if b.UpdatedAt != "" {
		metaItem["updatedAt"] = &types.AttributeValueMemberS{Value: b.UpdatedAt}
	}

	_, err := r.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:           aws.String(r.TableName),
					Item:                emailItem,
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
			{
				Put: &types.Put{
					TableName:           aws.String(r.TableName),
					Item:                phoneItem,
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
			{
				Put: &types.Put{
					TableName:           aws.String(r.TableName),
					Item:                apiKeyItem,
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
			{
				Put: &types.Put{
					TableName:           aws.String(r.TableName),
					Item:                metaItem,
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("transact create: %w", err)
	}

	return nil
}

func (r *BusinessRepository) GetByAPIKey(ctx context.Context, apiKey string) (*models.Business, error) {
	out, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "APIKEY#" + apiKey},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, fmt.Errorf("business not found")
	}

	businessPK := out.Items[0]["SK"].(*types.AttributeValueMemberS).Value

	// Obtener metadata del negocio
	return r.GetByPK(ctx, businessPK)
}

func (r *BusinessRepository) GetByPK(ctx context.Context, pk string) (*models.Business, error) {
	out, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, fmt.Errorf("business not found")
	}

	business := &models.Business{
		PK:        out.Item["PK"].(*types.AttributeValueMemberS).Value,
		SK:        out.Item["SK"].(*types.AttributeValueMemberS).Value,
		Name:      out.Item["name"].(*types.AttributeValueMemberS).Value,
		Email:     out.Item["email"].(*types.AttributeValueMemberS).Value,
		Phone:     out.Item["phone"].(*types.AttributeValueMemberS).Value,
		PlanID:    out.Item["planId"].(*types.AttributeValueMemberS).Value,
		APIKey:    out.Item["apiKey"].(*types.AttributeValueMemberS).Value,
		CreatedAt: out.Item["createdAt"].(*types.AttributeValueMemberS).Value,
	}

	if updatedAt, ok := out.Item["updatedAt"]; ok {
		business.UpdatedAt = updatedAt.(*types.AttributeValueMemberS).Value
	}

	return business, nil
}

func (r *BusinessRepository) UpdateAPIKey(ctx context.Context, businessPK, oldAPIKey, newAPIKey, updatedAt string) error {
	// Eliminar índice de API Key anterior
	deleteOldKeyItem := types.TransactWriteItem{
		Delete: &types.Delete{
			TableName: aws.String(r.TableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: "APIKEY#" + oldAPIKey},
				"SK": &types.AttributeValueMemberS{Value: businessPK},
			},
		},
	}

	// Crear nuevo índice de API Key
	newKeyItem := types.TransactWriteItem{
		Put: &types.Put{
			TableName: aws.String(r.TableName),
			Item: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: "APIKEY#" + newAPIKey},
				"SK": &types.AttributeValueMemberS{Value: businessPK},
			},
			ConditionExpression: aws.String("attribute_not_exists(PK)"),
		},
	}

	// Actualizar metadata con nueva API Key
	updateMetaItem := types.TransactWriteItem{
		Update: &types.Update{
			TableName: aws.String(r.TableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: businessPK},
				"SK": &types.AttributeValueMemberS{Value: "METADATA"},
			},
			UpdateExpression: aws.String("SET apiKey = :apiKey, updatedAt = :updatedAt"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":apiKey":    &types.AttributeValueMemberS{Value: newAPIKey},
				":updatedAt": &types.AttributeValueMemberS{Value: updatedAt},
			},
		},
	}

	_, err := r.Client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			deleteOldKeyItem,
			newKeyItem,
			updateMetaItem,
		},
	})

	if err != nil {
		return fmt.Errorf("transact update api key: %w", err)
	}

	return nil
}
