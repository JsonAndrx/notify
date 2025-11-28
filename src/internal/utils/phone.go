package utils

import (
	"regexp"
	"strings"
)

// ValidatePhoneNumber valida que el número de teléfono tenga un formato internacional válido
// Formato esperado: +[código país][número] (ej: +573001234567, +12025551234)
func ValidatePhoneNumber(phone string) bool {
	// Eliminar espacios en blanco
	phone = strings.TrimSpace(phone)

	// Debe empezar con +
	if !strings.HasPrefix(phone, "+") {
		return false
	}

	// Debe tener entre 10 y 15 dígitos (después del +)
	// Formato E.164: + seguido de 1-15 dígitos
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{9,14}$`)

	return phoneRegex.MatchString(phone)
}

// FormatPhoneNumber formatea el número eliminando espacios y caracteres especiales
func FormatPhoneNumber(phone string) string {
	// Eliminar espacios, guiones, paréntesis
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	return phone
}
