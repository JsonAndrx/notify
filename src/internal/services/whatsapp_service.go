package services

import (
	"context"
	"encoding/json"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"time"
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
// Ejemplo: {"1":"Juan Pérez","2":"Juan","3":"cita médica","4":"15 de dic","5":"+573001234567"}
func buildTwilioContentVariables(paramOrder []string, paramValues map[string]string) string {
	// Crear mapa con índices numéricos (1, 2, 3, ...)
	variables := make(map[string]string)
	for i, paramName := range paramOrder {
		// Twilio usa índices basados en 1 (no en 0)
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

	// TODO: Integración real con Twilio WhatsApp
	// Aquí deberías enviar a Twilio usando:
	//
	// import "github.com/twilio/twilio-go"
	// import twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	//
	// twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
	//     Username: os.Getenv("TWILIO_ACCOUNT_SID"),
	//     Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	// })
	//
	// params := &twilioApi.CreateMessageParams{}
	// params.SetTo("whatsapp:" + req.To)
	// params.SetFrom("whatsapp:" + os.Getenv("TWILIO_WHATSAPP_NUMBER"))
	// params.SetContentSid(template.ExternalID) // HXee830abc1d548d784f6963529c36b327
	// params.SetContentVariables(contentVariables) // {"1":"Juan Pérez","2":"Juan",...}
	//
	// message, err := twilioClient.Api.CreateMessage(params)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to send WhatsApp message: %v", err)
	// }
	// notificationID := *message.Sid

	// Por ahora simulamos el envío exitoso
	// En producción, aquí usarías el SID del mensaje de Twilio
	notificationID := fmt.Sprintf("WA_%d", time.Now().UnixNano())

	// Log para debugging (eliminar en producción o usar logger apropiado)
	fmt.Printf("WhatsApp simulation - To: %s, Template: %s, ContentSID: %s, Variables: %s\n",
		req.To, template.TemplateID, template.ExternalID, contentVariables)

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
