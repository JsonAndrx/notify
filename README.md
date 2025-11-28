# Notification Service - Backend API

Sistema de notificaciones con autenticación mediante API Keys y gestión de planes.

## 🚀 Características

- ✅ Registro de negocios con generación automática de API Keys
- ✅ Sistema de planes (FREE: 50 notificaciones / 30 días)
- ✅ Regeneración de API Keys
- ✅ Consulta de información de cuenta
- ✅ Consulta de uso del plan
- ✅ Envío de notificaciones con validación de límites
- ✅ Tracking automático de uso por período

## 📋 Estructura del Proyecto

```
src/
├── cmd/
│   ├── auth/register/          # Registro de negocios
│   ├── account/
│   │   ├── info/               # Info de cuenta
│   │   └── regenerate-key/     # Regenerar API Key
│   ├── plan/usage/             # Uso del plan
│   └── notifications/send/     # Enviar notificación
├── internal/
│   ├── models/                 # Modelos de datos
│   ├── repository/             # Acceso a DynamoDB
│   ├── services/               # Lógica de negocio
│   └── utils/                  # Utilidades (generación API Keys)
└── common/response/            # Respuestas HTTP
```

## 🛠️ Setup Inicial

### 1. Inicializar DynamoDB Local

```bash
# Ejecutar DynamoDB local en Docker
docker run -p 8000:8000 amazon/dynamodb-local
```

### 2. Crear la tabla

```bash
# Crear tabla NotificationService
aws dynamodb create-table \
  --table-name NotificationService \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --endpoint-url http://localhost:8000
```

### 3. Inicializar el Plan FREE

```bash
chmod +x scripts/init_free_plan.sh
./scripts/init_free_plan.sh
```

### 4. Inicializar Templates de WhatsApp

```bash
chmod +x scripts/init_templates.sh
./scripts/init_templates.sh
```

### 5. Compilar y Desplegar

```bash
cd infrastructure
sam build
sam local start-api --env-vars env.json
```

## 📡 API Endpoints

### 1. Registro de Negocio

**POST** `/v1/business/register`

```json
{
  "name": "Mi Empresa",
  "email": "contacto@miempresa.com",
  "phone": "+1234567890",
  "plan_id": "FREE"
}
```

**Respuesta:**
```json
{
  "id_business": "uuid",
  "name": "Mi Empresa",
  "email": "contacto@miempresa.com",
  "phone": "+1234567890",
  "plan_id": "FREE",
  "api_key": "nfy_..."
}
```

### 2. Regenerar API Key

**POST** `/v1/account/regenerate-key`

Por seguridad, requiere email, phone y API Key actual para verificar identidad.

```json
{
  "email": "contacto@miempresa.com",
  "phone": "+1234567890",
  "current_api_key": "nfy_..."
}
```

**Respuesta:**
```json
{
  "api_key": "nfy_new_key..."
}
```

**Códigos de Error:**
- `401`: Credenciales inválidas (email, phone o API key no coinciden)
- `400`: Request body inválido

### 3. Información de la Cuenta

**GET** `/v1/account/info`

**Headers:**
```
X-API-Key: nfy_...
```

**Respuesta:**
```json
{
  "id_business": "uuid",
  "name": "Mi Empresa",
  "email": "contacto@miempresa.com",
  "phone": "+1234567890",
  "plan_id": "FREE",
  "created_at": "2025-11-26T10:00:00Z"
}
```

### 4. Uso del Plan

**GET** `/v1/plan/usage`

**Headers:**
```
X-API-Key: nfy_...
```

**Respuesta:**
```json
{
  "business_id": "uuid",
  "plan_id": "FREE",
  "plan_name": "Free Plan",
  "notification_limit": 50,
  "notification_count": 10,
  "notification_left": 40,
  "period_start": "2025-11-01T00:00:00Z",
  "period_end": "2025-12-01T00:00:00Z",
  "period_days": 30
}
```

### 5. Enviar Notificación WhatsApp (con Template)

**POST** `/v1/notifications/whatsapp`

**Headers:**
```
X-API-Key: nfy_...
```

**Body:**
```json
{
  "to": "+1234567890",
  "template_id": "whatsapp-verification-code",
  "parameters": {
    "name": "Juan",
    "code": "123456"
  }
}
```

**Respuesta:**
```json
{
  "success": true,
  "notification_id": "WA_...",
  "template_used": "Código de Verificación",
  "notification_count": 11,
  "notification_left": 39
}
```

**Códigos de Error:**
- `401`: API Key inválida
- `429`: Límite de notificaciones alcanzado
- `404`: Template no encontrado o inactivo
- `400`: Parámetros inválidos o faltantes

### 6. Enviar SMS

**POST** `/v1/notifications/sms`

**Headers:**
```
X-API-Key: nfy_...
```

**Body:**
```json
{
  "to": "+1234567890",
  "message": "Tu código de verificación es: 123456"
}
```

**Respuesta:**
```json
{
  "success": true,
  "notification_id": "SMS_...",
  "notification_count": 12,
  "notification_left": 38
}
```

### 7. Enviar Email

**POST** `/v1/notifications/email`

**Headers:**
```
X-API-Key: nfy_...
```

**Body:**
```json
{
  "to": "usuario@example.com",
  "subject": "Bienvenido a nuestro servicio",
  "body": "<h1>Hola!</h1><p>Gracias por registrarte</p>",
  "html": true
}
```

**Respuesta:**
```json
{
  "success": true,
  "notification_id": "EMAIL_...",
  "notification_count": 13,
  "notification_left": 37
}
```

## 🗃️ Estructura de Datos en DynamoDB

### Business
```
PK: BUSINESS#{uuid}
SK: METADATA
name, email, phone, planId, apiKey, createdAt, updatedAt
```

### Índices de Búsqueda
```
PK: EMAIL#{email}
SK: BUSINESS#{uuid}

PK: PHONE#{phone}
SK: BUSINESS#{uuid}

PK: APIKEY#{apiKey}
SK: BUSINESS#{uuid}
```

### Plan
```
PK: PLAN#{planId}
SK: METADATA
name, notificationLimit, periodDays, price, description, active, createdAt
```

### Usage
```
PK: BUSINESS#{uuid}
SK: USAGE#{date}
businessId, planId, notificationCount, periodStart, periodEnd, createdAt, updatedAt
```

### Template
```
PK: TEMPLATE#{templateId}
SK: METADATA
templateId, name, type, provider, externalId, parameters[], parameterCount, description, active, createdAt, updatedAt
```

## 📋 Plantillas de WhatsApp

El sistema soporta plantillas de WhatsApp con validación de parámetros. Las plantillas se configuran con:

- **ID único**: Identificador de la plantilla
- **Tipo**: whatsapp, sms, email
- **Provider**: twilio, sendgrid, etc.
- **External ID**: ID de la plantilla en el proveedor externo
- **Parámetros**: Lista de parámetros requeridos
- **Validación**: Automática de parámetros faltantes o extra

### Plantillas Incluidas

#### 1. Código de Verificación
- **ID**: `whatsapp-verification-code`
- **Parámetros**: `name`, `code`
- **Uso**: Envío de códigos de autenticación

```json
{
  "to": "+1234567890",
  "template_id": "whatsapp-verification-code",
  "parameters": {
    "name": "Juan",
    "code": "123456"
  }
}
```

#### 2. Confirmación de Pedido
- **ID**: `whatsapp-order-confirmation`
- **Parámetros**: `customer_name`, `order_number`, `total_amount`
- **Uso**: Confirmación de compras

#### 3. Recordatorio de Cita
- **ID**: `whatsapp-appointment-reminder`
- **Parámetros**: `name`, `date`, `time`, `location`
- **Uso**: Recordatorios de citas médicas, reuniones, etc.

#### 4. Mensaje de Bienvenida
- **ID**: `whatsapp-welcome`
- **Parámetros**: `name`
- **Uso**: Onboarding de nuevos usuarios

## 🔒 Seguridad

- Las API Keys se generan con `crypto/rand` (32 bytes)
- Formato: `nfy_` + base64 URL-safe
- Cada negocio tiene un único API Key activo
- Las API Keys se validan en cada request

## 📝 Notas Importantes

1. **Plan FREE por defecto**: Si no se especifica `plan_id` en el registro, se asigna automáticamente el plan FREE
2. **Períodos automáticos**: El sistema crea automáticamente nuevos períodos de uso cuando expira el actual
3. **Límites**: El plan FREE permite 50 notificaciones cada 30 días
4. **Renovación**: Al finalizar un período, el contador se reinicia automáticamente

## 🚧 Próximas Mejoras

- [ ] Implementar envío real de WhatsApp, SMS y Email
- [ ] Agregar webhooks para notificaciones
- [ ] Implementar más planes (BASIC, PRO, ENTERPRISE)
- [ ] Agregar templates de mensajes
- [ ] Historial de notificaciones enviadas
- [ ] Dashboard de estadísticas

## 🧪 Testing

```bash
# 1. Registrar un negocio
curl -X POST http://localhost:3000/v1/business/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Business",
    "email": "test@example.com",
    "phone": "+1234567890"
  }'

# Respuesta (guardar el api_key):
# {"id_business":"xxx","name":"Test Business","email":"test@example.com","phone":"+1234567890","plan_id":"FREE","api_key":"nfy_..."}

# 2. Consultar información de la cuenta
curl -X GET http://localhost:3000/v1/account/info \
  -H "X-API-Key: nfy_..."

# 3. Consultar uso del plan
curl -X GET http://localhost:3000/v1/plan/usage \
  -H "X-API-Key: nfy_..."

# 4. Enviar notificación
curl -X POST http://localhost:3000/v1/notifications/send \
  -H "X-API-Key: nfy_..." \
  -H "Content-Type: application/json" \
  -d '{
    "to": "+9876543210",
    "message": "Test notification",
    "type": "whatsapp"
  }'

# 5. Regenerar API Key (requiere verificación)
curl -X POST http://localhost:3000/v1/account/regenerate-key \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "phone": "+1234567890",
    "current_api_key": "nfy_old_key..."
  }'

# Respuesta: nueva API Key
# {"api_key":"nfy_new_key..."}
```
