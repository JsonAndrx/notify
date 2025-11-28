package services

import (
	"os"

	"github.com/twilio/twilio-go"
)

// GetTwilioClient crea y retorna un cliente de Twilio configurado
func GetTwilioClient() *twilio.RestClient {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	if accountSid == "" || authToken == "" {
		// Si no hay credenciales, retornar nil para modo simulación
		return nil
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return client
}

// GetTwilioWhatsAppNumber obtiene el número de WhatsApp de Twilio desde variables de entorno
func GetTwilioWhatsAppNumber() string {
	return os.Getenv("TWILIO_WHATSAPP_NUMBER")
}
