package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// GenerateAPIKey genera una API Key única y segura
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convertir a base64 y limpiar caracteres problemáticos
	apiKey := base64.URLEncoding.EncodeToString(b)
	apiKey = strings.TrimRight(apiKey, "=")

	// Formato: nfy_<key>
	return "nfy_" + apiKey, nil
}

// ValidateAPIKeyFormat valida que la API Key tenga el formato correcto
func ValidateAPIKeyFormat(apiKey string) bool {
	if !strings.HasPrefix(apiKey, "nfy_") {
		return false
	}

	if len(apiKey) < 10 {
		return false
	}

	return true
}

// ExtractAPIKeyFromRequest busca el API Key en los headers de un request
// Busca en múltiples formatos: X-API-Key, x-api-key, y case-insensitive
func ExtractAPIKeyFromRequest(request events.APIGatewayProxyRequest) string {
	// Intentar en Headers (case-sensitive)
	if apiKey := request.Headers["X-API-Key"]; apiKey != "" {
		return apiKey
	}
	if apiKey := request.Headers["x-api-key"]; apiKey != "" {
		return apiKey
	}

	// Intentar búsqueda case-insensitive en todos los headers
	for key, value := range request.Headers {
		if strings.ToLower(key) == "x-api-key" {
			return value
		}
	}

	return ""
}
