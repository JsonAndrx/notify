#!/bin/bash

# Script para inicializar templates de WhatsApp en DynamoDB

set -e

TABLE_NAME="NotificationService"
ENDPOINT="http://localhost:8000"

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "🔧 Inicializando Templates de Notificaciones"
echo "============================================="
echo ""

# Template 1: Recordatorio General (WhatsApp)
echo -n "Creando template: recordatorio_general... "
aws dynamodb put-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
    --region us-east-1 \
    --item '{
        "PK": {"S": "TEMPLATE#recordatorio_general"},
        "SK": {"S": "METADATA"},
        "templateId": {"S": "recordatorio_general"},
        "name": {"S": "Recordatorio General"},
        "type": {"S": "whatsapp"},
        "provider": {"S": "twilio"},
        "externalId": {"S": "HXee830abc1d548d784f6963529c36b327"},
        "parameters": {"L": [
            {"S": "header_name"},
            {"S": "body_name"},
            {"S": "service"},
            {"S": "date"},
            {"S": "contact"}
        ]},
        "parameterCount": {"N": "5"},
        "description": {"S": "Header={{1}}: nombre | Body: 👋 ¡Hola {{1}}=nombre!, te recordamos tu {{2}}=servicio programado(a) para {{3}}=fecha. ⏰ contáctanos a {{4}}=contacto, Si necesitas reprogramar."},
        "active": {"BOOL": true},
        "createdAt": {"S": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}
    }' > /dev/null

echo -e "${GREEN}✓${NC}"

# Template 2: Actualización de Estado (WhatsApp)
echo -n "Creando template: actualizacion_estado... "
aws dynamodb put-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
    --region us-east-1 \
    --item '{
        "PK": {"S": "TEMPLATE#actualizacion_estado"},
        "SK": {"S": "METADATA"},
        "templateId": {"S": "actualizacion_estado"},
        "name": {"S": "Actualización de Estado"},
        "type": {"S": "whatsapp"},
        "provider": {"S": "twilio"},
        "externalId": {"S": "HX78b9865cd3ad336a37d17fa63b7f0c26"},
        "parameters": {"L": [
            {"S": "header_name"},
            {"S": "body_name"},
            {"S": "service"},
            {"S": "status"},
            {"S": "company"}
        ]},
        "parameterCount": {"N": "5"},
        "description": {"S": "Header={{1}}: nombre | Body: 👋 ¡Hola {{1}}=nombre!, tenemos una actualización sobre tu {{2}}=servicio. 📍 Estado actual: {{3}}=estado. Somos {{4}}=empresa, te mantendremos informado(a)."},
        "active": {"BOOL": true},
        "createdAt": {"S": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}
    }' > /dev/null

echo -e "${GREEN}✓${NC}"

# Template 3: Confirmación General (WhatsApp)
echo -n "Creando template: confirmacion_general... "
aws dynamodb put-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
    --region us-east-1 \
    --item '{
        "PK": {"S": "TEMPLATE#confirmacion_general"},
        "SK": {"S": "METADATA"},
        "templateId": {"S": "confirmacion_general"},
        "name": {"S": "Confirmación General"},
        "type": {"S": "whatsapp"},
        "provider": {"S": "twilio"},
        "externalId": {"S": "HXcbd993c4737e3c2a9b0f9a7332f01b33"},
        "parameters": {"L": [
            {"S": "header_name"},
            {"S": "body_name"},
            {"S": "service"},
            {"S": "detail"},
            {"S": "company"}
        ]},
        "parameterCount": {"N": "5"},
        "description": {"S": "Header={{1}}: nombre | Body: 👋 ¡Hola {{1}}=nombre!, tu {{2}}=servicio ha sido confirmado(a). 📌 Detalle: {{3}}=detalle. Somos {{4}}=empresa gracias por confiar en nosotros."},
        "active": {"BOOL": true},
        "createdAt": {"S": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}
    }' > /dev/null

echo -e "${GREEN}✓${NC}"

echo ""
echo -e "${GREEN}✓ Templates creados exitosamente${NC}"
echo ""
echo "Templates disponibles:"
echo "  1. recordatorio_general"
echo "     - header_name: {{1}} del header"
echo "     - body_name: {{1}} del body"
echo "     - service: {{2}} del body"
echo "     - date: {{3}} del body"
echo "     - contact: {{4}} del body"
echo ""
echo "  2. actualizacion_estado"
echo "     - header_name: {{1}} del header"
echo "     - body_name: {{1}} del body"
echo "     - service: {{2}} del body"
echo "     - status: {{3}} del body"
echo "     - company: {{4}} del body"
echo ""
echo "  3. confirmacion_general"
echo "     - header_name: {{1}} del header"
echo "     - body_name: {{1}} del body"
echo "     - service: {{2}} del body"
echo "     - detail: {{3}} del body"
echo "     - company: {{4}} del body"
echo ""
echo "Ejemplo de uso:"
echo 'curl -X POST http://localhost:3000/v1/notifications/whatsapp \'
echo '  -H "X-API-Key: nfy_..." \'
echo '  -H "Content-Type: application/json" \'
echo '  -d '"'"'{
    "to": "+573001234567",
    "template_id": "recordatorio_general",
    "parameters": {
      "header_name": "Juan Pérez",
      "body_name": "Juan",
      "service": "cita médica",
      "date": "15 de diciembre a las 10:00 AM",
      "contact": "+573001234567"
    }
  }'"'"
