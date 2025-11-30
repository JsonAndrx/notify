package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"notify-backend/internal/db"
	"notify-backend/internal/repository"
	"os"
	"strings"
	"time"
)

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	HTML    bool   `json:"html"` // Si el body es HTML o texto plano
}

type SendEmailResponse struct {
	Success           bool   `json:"success"`
	NotificationID    string `json:"notification_id"`
	NotificationCount int    `json:"notification_count"`
	NotificationLeft  int    `json:"notification_left"`
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// sendEmailSMTP envía un email usando SMTP (Gmail)
func sendEmailSMTP(host, port, username, password, to, subject, body string, isHTML bool, fromName string) error {
	// Configurar autenticación
	auth := smtp.PlainAuth("", username, password, host)

	// Construir el mensaje con formato correcto
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	// Usar formato RFC 5322 correcto
	from := username
	if fromName != "" {
		from = fmt.Sprintf("%s <%s>", fromName, username)
	}

	// Construir headers del email
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("%s; charset=UTF-8", contentType)
	headers["Content-Transfer-Encoding"] = "quoted-printable"

	// Construir el mensaje completo
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	msg := []byte(message)

	// Para Gmail, necesitamos usar STARTTLS
	if host == "smtp.gmail.com" {
		return sendWithSTARTTLS(host, port, auth, username, []string{to}, msg)
	}

	// Para otros servidores SMTP estándar
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, username, []string{to}, msg)
}

// sendWithSTARTTLS envía email usando STARTTLS (requerido por Gmail)
func sendWithSTARTTLS(host, port string, auth smtp.Auth, from string, to []string, msg []byte) error {
	addr := fmt.Sprintf("%s:%s", host, port)

	// Conectar al servidor
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	if err = conn.Hello("localhost"); err != nil {
		return fmt.Errorf("failed EHLO: %v", err)
	}

	// Iniciar TLS
	fmt.Printf("   → Starting TLS...\n")
	tlsConfig := &tls.Config{
		ServerName: host,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed STARTTLS: %v", err)
	}

	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("failed authentication: %v", err)
	}

	if err = conn.Mail(from); err != nil {
		return fmt.Errorf("failed MAIL: %v", err)
	}

	for _, recipient := range to {
		if err = conn.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed RCPT: %v", err)
		}
	}

	writer, err := conn.Data()
	if err != nil {
		return fmt.Errorf("failed DATA: %v", err)
	}
	defer writer.Close()

	if _, err = writer.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}
	fmt.Printf("   ✓ Message data sent\n")

	return nil
}

func SendEmailService(apiKey string, req SendEmailRequest) (*SendEmailResponse, error) {
	client, _ := db.NewDynamoClient()
	businessRepo := repository.NewBusinessRepository(client, "NotificationService")
	planRepo := repository.NewPlanRepository(client, "NotificationService")
	usageRepo := repository.NewUsageRepository(client, "NotificationService")
	ctx := context.TODO()

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

	// Validar email del destinatario
	if !strings.Contains(req.To, "@") || !strings.Contains(req.To, ".") {
		return nil, fmt.Errorf("invalid email address")
	}

	// Validaciones adicionales
	if req.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}
	if req.Body == "" {
		return nil, fmt.Errorf("body is required")
	}

	// Enviar email por SMTP
	var notificationID string
	smtpHost := getEnv("SMTP_HOST", "")
	smtpPort := getEnv("SMTP_PORT", "")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPass := getEnv("SMTP_PASSWORD", "")


	if smtpUser != "" && smtpPass != "" {
		// Envío real por SMTP
		err := sendEmailSMTP(smtpHost, smtpPort, smtpUser, smtpPass, req.To, req.Subject, req.Body, req.HTML, business.Name)
		if err != nil {
			fmt.Printf("❌ Failed to send email: %v\n", err)
			return nil, fmt.Errorf("failed to send notification")
		}
		notificationID = fmt.Sprintf("EMAIL_%d", time.Now().UnixNano())
		fmt.Printf("✅ Email sent successfully!\n")
		fmt.Printf("   Message ID: %s\n", notificationID)
	} else {
		return nil, fmt.Errorf("email service not configured")
	}

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

	return &SendEmailResponse{
		Success:           true,
		NotificationID:    notificationID,
		NotificationCount: newCount,
		NotificationLeft:  notificationLeft,
	}, nil
}
