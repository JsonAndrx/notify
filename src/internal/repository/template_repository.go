package repository

import (
	"context"
	"fmt"

	"notify-backend/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TemplateRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewTemplateRepository(client *dynamodb.Client, tableName string) *TemplateRepository {
	return &TemplateRepository{
		Client:    client,
		TableName: tableName,
	}
}

// GetByID obtiene una plantilla por su ID
func (r *TemplateRepository) GetByID(ctx context.Context, templateID string) (*models.Template, error) {
	out, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TEMPLATE#" + templateID},
			"SK": &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, fmt.Errorf("template not found")
	}

	var template models.Template
	err = attributevalue.UnmarshalMap(out.Item, &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// GetByTypeAndExternalID obtiene una plantilla por tipo y external ID (para búsquedas específicas)
func (r *TemplateRepository) GetByTypeAndExternalID(ctx context.Context, templateType, externalID string) (*models.Template, error) {
	// Esto requeriría un GSI, por ahora hacemos un scan (no recomendado en producción)
	out, err := r.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.TableName),
		FilterExpression: aws.String("begins_with(PK, :pk) AND #type = :type AND externalId = :externalId"),
		ExpressionAttributeNames: map[string]string{
			"#type": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":         &types.AttributeValueMemberS{Value: "TEMPLATE#"},
			":type":       &types.AttributeValueMemberS{Value: templateType},
			":externalId": &types.AttributeValueMemberS{Value: externalID},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, fmt.Errorf("template not found")
	}

	var template models.Template
	err = attributevalue.UnmarshalMap(out.Items[0], &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// ListByType lista todas las plantillas de un tipo específico
func (r *TemplateRepository) ListByType(ctx context.Context, templateType string) ([]*models.Template, error) {
	out, err := r.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.TableName),
		FilterExpression: aws.String("begins_with(PK, :pk) AND #type = :type AND active = :active"),
		ExpressionAttributeNames: map[string]string{
			"#type": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: "TEMPLATE#"},
			":type":   &types.AttributeValueMemberS{Value: templateType},
			":active": &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		return nil, err
	}

	templates := make([]*models.Template, 0, len(out.Items))
	for _, item := range out.Items {
		var template models.Template
		err = attributevalue.UnmarshalMap(item, &template)
		if err != nil {
			continue
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// Create crea una nueva plantilla
func (r *TemplateRepository) Create(ctx context.Context, template *models.Template) error {
	item, err := attributevalue.MarshalMap(template)
	if err != nil {
		return err
	}

	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	})

	return err
}

// ValidateTemplateParameters valida que los parámetros proporcionados coincidan con la plantilla
func (r *TemplateRepository) ValidateTemplateParameters(template *models.Template, providedParams map[string]string) *models.TemplateValidation {
	validation := &models.TemplateValidation{
		Valid:          true,
		MissingParams:  []string{},
		ExtraParams:    []string{},
		ParameterCount: len(providedParams),
		ExpectedCount:  template.ParameterCount,
	}

	// Verificar parámetros faltantes
	for _, requiredParam := range template.Parameters {
		if _, exists := providedParams[requiredParam]; !exists {
			validation.MissingParams = append(validation.MissingParams, requiredParam)
			validation.Valid = false
		}
	}

	// Verificar parámetros extra
	requiredMap := make(map[string]bool)
	for _, param := range template.Parameters {
		requiredMap[param] = true
	}

	for providedParam := range providedParams {
		if !requiredMap[providedParam] {
			validation.ExtraParams = append(validation.ExtraParams, providedParam)
		}
	}

	return validation
}
