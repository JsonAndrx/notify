package services

import (
	"context"
	"fmt"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"notify-backend/internal/utils"
	"regexp"
	"strings"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SendSMSRequest struct {
	To         string            `json:"to"`
	TemplateID string            `json:"template_id"`
	Parameters map[string]string `json:"parameters"`
}

type SendSMSResponse struct {
	Success           bool   `json:"success"`
	NotificationID    string `json:"notification_id"`
	TemplateUsed      string `json:"template_used"`
	NotificationCount int    `json:"notification_count"`
	NotificationLeft  int    `json:"notification_left"`
}

// buildSMSMessage construye el mensaje SMS a partir del template y parámetros
func buildSMSMessage(template string, params map[string]string) string {
	message := template
	for key, value := range params {
		placeholder := "$" + key
		message = strings.ReplaceAll(message, placeholder, value)
	}
	return message
}

// validateVerificationCode valida que el código sea numérico y tenga entre 4 y 6 dígitos
func validateVerificationCode(code string) bool {
	matched, _ := regexp.MatchString(`^\d{4,6}$`, code)
	return matched
}

func SendSMSService(apiKey string, req SendSMSRequest) (*SendSMSResponse, error) {
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

	// Validar que la plantilla existe y es de tipo sms
	template, err := templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("invalid template")
	}

	if template.Type != "sms" {
		return nil, fmt.Errorf("invalid template type")
	}

	if !template.Active {
		return nil, fmt.Errorf("template not available")
	}

	// Validar parámetros de la plantilla
	validation := templateRepo.ValidateTemplateParameters(template, req.Parameters)
	if !validation.Valid {
		if len(validation.MissingParams) > 0 {
			return nil, fmt.Errorf("missing required parameters")
		}
	}

	// Validaciones específicas para template de verificación
	if template.TemplateID == "sms_verification_code" {
		code, hasCode := req.Parameters["codigo"]
		if !hasCode || !validateVerificationCode(code) {
			return nil, fmt.Errorf("invalid verification code format")
		}
	}

	// Construir mensaje desde template
	// Agregar el nombre de la empresa automáticamente
	req.Parameters["empresa"] = business.Name
	message := buildSMSMessage(template.Description, req.Parameters)

	// Validar longitud del mensaje final
	if len(message) > 1600 {
		return nil, fmt.Errorf("message too long")
	}

	// Obtener cliente de Twilio
	twilioClient := GetTwilioClient()
	twilioPhoneNumber := GetTwilioPhoneNumber()

	// Verificar que Twilio esté configurado
	if twilioClient == nil || twilioPhoneNumber == "" {
		return nil, fmt.Errorf("service temporarily unavailable")
	}

	// Enviar SMS a través de Twilio
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(req.To)
	params.SetFrom(twilioPhoneNumber)
	params.SetBody(message)

	twilioMessage, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		// Log interno del error real para debugging
		fmt.Printf("Twilio SMS error: %v\n", err)
		return nil, fmt.Errorf("failed to send notification")
	}

	notificationID := *twilioMessage.Sid
	fmt.Printf("SMS sent - MessageSID: %s, To: %s, Template: %s\n", notificationID, req.To, template.TemplateID)

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

	return &SendSMSResponse{
		Success:           true,
		NotificationID:    notificationID,
		TemplateUsed:      template.Name,
		NotificationCount: newCount,
		NotificationLeft:  notificationLeft,
	}, nil
}
