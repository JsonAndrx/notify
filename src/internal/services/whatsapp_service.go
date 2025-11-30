package services

import (
	"context"
	"encoding/json"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"notify-backend/internal/utils"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SendWhatsAppRequest struct {
	To         string            `json:"to"`
	TemplateID string            `json:"template_id"`
	Parameters map[string]string `json:"parameters"`
}

type SendWhatsAppResponse struct {
	Success           bool   `json:"success"`
	NotificationID    string `json:"notification_id"`
	TemplateUsed      string `json:"template_used"`
	NotificationCount int    `json:"notification_count"`
	NotificationLeft  int    `json:"notification_left"`
}

// buildTwilioContentVariables convierte parámetros nombrados a formato JSON de Twilio
// En WhatsApp, cuando el template tiene header y body que usan {{1}}, ambos comparten el mismo valor
// Ejemplo: template con 3 params ["name", "service", "date"] genera {"1":"Juan","2":"cita","3":"10:00"}
// Donde {{1}} del header y {{1}} del body serán "Juan"
func buildTwilioContentVariables(paramOrder []string, paramValues map[string]string) string {
	variables := make(map[string]string)

	// Mapear directamente en orden: 1, 2, 3, ...
	for i, paramName := range paramOrder {
		variables[fmt.Sprintf("%d", i+1)] = paramValues[paramName]
	}

	jsonBytes, _ := json.Marshal(variables)
	return string(jsonBytes)
}

func SendWhatsAppService(apiKey string, req SendWhatsAppRequest) (*SendWhatsAppResponse, error) {
	client, _ := db.NewDynamoClient()
	businessRepo := repository.NewBusinessRepository(client, "NotificationService")
	planRepo := repository.NewPlanRepository(client, "NotificationService")
	usageRepo := repository.NewUsageRepository(client, "NotificationService")
	templateRepo := repository.NewTemplateRepository(client, "NotificationService")
	ctx := context.TODO()

	// Validar formato del número de teléfono
	formattedPhone := utils.FormatPhoneNumber(req.To)
	if !utils.ValidatePhoneNumber(formattedPhone) {
		return nil, fmt.Errorf("invalid phone number format")
	}
	req.To = formattedPhone

	// Buscar negocio por API Key
	business, err := businessRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("authentication failed")
	}

	businessID := business.PK[9:] // Remover "BUSINESS#"

	// Obtener plan
	plan, err := planRepo.GetByID(ctx, business.PlanID)
	if err != nil {
		return nil, fmt.Errorf("service unavailable")
	}

	// Verificar o crear período de uso
	usage, err := usageRepo.CheckAndCreateNewPeriod(ctx, businessID, business.PlanID, plan.PeriodDays)
	if err != nil {
		return nil, fmt.Errorf("service unavailable")
	}

	// Verificar límite de notificaciones
	if usage.NotificationCount >= plan.NotificationLimit {
		return nil, fmt.Errorf("notification limit reached")
	}

	// Validar que la plantilla existe y es de tipo whatsapp
	template, err := templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}

	if template.Type != "whatsapp" {
		return nil, fmt.Errorf("invalid template")
	}

	if !template.Active {
		return nil, fmt.Errorf("template not available")
	}

	// Validar parámetros de la plantilla
	validation := templateRepo.ValidateTemplateParameters(template, req.Parameters)
	if !validation.Valid {
		if len(validation.MissingParams) > 0 {
			return nil, fmt.Errorf("invalid template parameters")
		}
	}

	// Convertir parámetros nombrados a formato Twilio
	// Twilio espera variables en formato: {"1":"valor1","2":"valor2","3":"valor3",...}
	contentVariables := buildTwilioContentVariables(template.Parameters, req.Parameters)

	// Obtener cliente de Twilio
	twilioClient := GetTwilioClient()
	twilioWhatsAppNumber := GetTwilioWhatsAppNumber()

	// Verificar que Twilio esté configurado
	if twilioClient == nil || twilioWhatsAppNumber == "" {
		return nil, fmt.Errorf("service temporarily unavailable")
	}

	// Enviar mensaje a través de Twilio WhatsApp
	params := &twilioApi.CreateMessageParams{}
	params.SetTo("whatsapp:" + req.To)
	params.SetFrom("whatsapp:" + twilioWhatsAppNumber)
	params.SetContentSid(template.ExternalID)
	params.SetContentVariables(contentVariables)

	message, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		// Log interno del error real para debugging
		fmt.Printf("Twilio WhatsApp error: %v\n", err)
		return nil, fmt.Errorf("failed to send notification")
	}

	notificationID := *message.Sid
	fmt.Printf("WhatsApp sent - MessageSID: %s, To: %s, Template: %s\n",
		notificationID, req.To, template.TemplateID)

	// Incrementar contador de uso
	err = usageRepo.IncrementUsage(ctx, businessID, usage.SK)
	if err != nil {
		// Log interno para debugging
		fmt.Printf("Failed to increment usage: %v\n", err)
		// No fallar la request si el mensaje ya fue enviado
	}

	// Calcular notificaciones restantes
	newCount := usage.NotificationCount + 1
	notificationLeft := plan.NotificationLimit - newCount
	if notificationLeft < 0 {
		notificationLeft = 0
	}

	return &SendWhatsAppResponse{
		Success:           true,
		NotificationID:    notificationID,
		TemplateUsed:      template.Name,
		NotificationCount: newCount,
		NotificationLeft:  notificationLeft,
	}, nil
}
