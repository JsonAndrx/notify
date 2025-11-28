package services

import (
	"context"
	"encoding/json"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"time"

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

func buildTwilioContentVariables(paramOrder []string, paramValues map[string]string) string {
	variables := make(map[string]string)

	idx := 1
	skipHeaderName := false

	// Verificar si tenemos tanto header_name como body_name
	hasHeaderName := false
	hasBodyName := false
	for _, paramName := range paramOrder {
		if paramName == "header_name" {
			hasHeaderName = true
		}
		if paramName == "body_name" {
			hasBodyName = true
		}
	}

	// Si tenemos ambos, skip header_name porque body_name será el {{1}} compartido
	skipHeaderName = hasHeaderName && hasBodyName

	for _, paramName := range paramOrder {
		// Saltar header_name si tenemos body_name
		if paramName == "header_name" && skipHeaderName {
			continue
		}

		variables[fmt.Sprintf("%d", idx)] = paramValues[paramName]
		idx++
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

	// Buscar negocio por API Key
	business, err := businessRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}

	businessID := business.PK[9:] // Remover "BUSINESS#"

	// Obtener plan
	plan, err := planRepo.GetByID(ctx, business.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan not found")
	}

	// Verificar o crear período de uso
	usage, err := usageRepo.CheckAndCreateNewPeriod(ctx, businessID, business.PlanID, plan.PeriodDays)
	if err != nil {
		return nil, fmt.Errorf("failed to check usage: %v", err)
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
		return nil, fmt.Errorf("invalid template type. Expected whatsapp, got %s", template.Type)
	}

	if !template.Active {
		return nil, fmt.Errorf("template is not active")
	}

	// Validar parámetros de la plantilla
	validation := templateRepo.ValidateTemplateParameters(template, req.Parameters)
	if !validation.Valid {
		if len(validation.MissingParams) > 0 {
			return nil, fmt.Errorf("missing required parameters: %v", validation.MissingParams)
		}
	}

	// Convertir parámetros nombrados a formato Twilio
	// Twilio espera variables en formato: {"1":"valor1","2":"valor2","3":"valor3",...}
	contentVariables := buildTwilioContentVariables(template.Parameters, req.Parameters)

	var notificationID string

	// Obtener cliente de Twilio
	twilioClient := GetTwilioClient()
	twilioWhatsAppNumber := GetTwilioWhatsAppNumber()

	if twilioClient != nil && twilioWhatsAppNumber != "" {
		// Integración real con Twilio WhatsApp
		params := &twilioApi.CreateMessageParams{}
		params.SetTo("whatsapp:" + req.To)
		params.SetFrom("whatsapp:" + twilioWhatsAppNumber)
		params.SetContentSid(template.ExternalID)
		params.SetContentVariables(contentVariables)

		message, err := twilioClient.Api.CreateMessage(params)
		if err != nil {
			return nil, fmt.Errorf("failed to send WhatsApp message: %v", err)
		}

		notificationID = *message.Sid
		fmt.Printf("WhatsApp sent - MessageSID: %s, To: %s, Template: %s\n",
			notificationID, req.To, template.TemplateID)
	} else {
		// Modo simulación (sin credenciales de Twilio)
		notificationID = fmt.Sprintf("WA_SIM_%d", time.Now().UnixNano())
		fmt.Printf("WhatsApp SIMULATION - To: %s, Template: %s, ContentSID: %s, Variables: %s\n",
			req.To, template.TemplateID, template.ExternalID, contentVariables)
		fmt.Println("💡 Tip: Configure TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN y TWILIO_WHATSAPP_NUMBER para envíos reales")
	}

	// Incrementar contador de uso
	err = usageRepo.IncrementUsage(ctx, businessID, usage.SK)
	if err != nil {
		return nil, fmt.Errorf("failed to increment usage: %v", err)
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
