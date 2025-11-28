package models

type Template struct {
	PK             string   `dynamodbav:"PK"`             // TEMPLATE#{templateId}
	SK             string   `dynamodbav:"SK"`             // METADATA
	TemplateID     string   `dynamodbav:"templateId"`     // ID único de la plantilla
	Name           string   `dynamodbav:"name"`           // Nombre descriptivo
	Type           string   `dynamodbav:"type"`           // whatsapp, sms, email
	Provider       string   `dynamodbav:"provider"`       // twilio, sendgrid, etc.
	ExternalID     string   `dynamodbav:"externalId"`     // ID de la plantilla en el proveedor (ej: Twilio template SID)
	Parameters     []string `dynamodbav:"parameters"`     // Lista de parámetros requeridos ["name", "code", "date"]
	ParameterCount int      `dynamodbav:"parameterCount"` // Número de parámetros
	Description    string   `dynamodbav:"description"`    // Descripción de la plantilla
	Active         bool     `dynamodbav:"active"`         // Si está activa o no
	CreatedAt      string   `dynamodbav:"createdAt"`
	UpdatedAt      string   `dynamodbav:"updatedAt,omitempty"`
}

// TemplateValidation representa la validación de una plantilla
type TemplateValidation struct {
	Valid          bool
	MissingParams  []string
	ExtraParams    []string
	ParameterCount int
	ExpectedCount  int
}
