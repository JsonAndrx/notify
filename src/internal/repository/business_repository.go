package repository

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "notify-backend/internal/models"
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
        TableName: aws.String(r.TableName),
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
        TableName: aws.String(r.TableName),
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

    metaItem := map[string]types.AttributeValue{
        "PK":        &types.AttributeValueMemberS{Value: b.PK},
        "SK":        &types.AttributeValueMemberS{Value: "METADATA"},
        "name":      &types.AttributeValueMemberS{Value: b.Name},
        "email":     &types.AttributeValueMemberS{Value: b.Email},
        "phone":     &types.AttributeValueMemberS{Value: b.Phone},
        "planId":    &types.AttributeValueMemberS{Value: b.PlanID},
        "createdAt": &types.AttributeValueMemberS{Value: b.CreatedAt},
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
