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
        "externalId": {"S": "HX334eac2cb8f3264c1cd104107bd39584"},
        "parameters": {"L": [
            {"S": "company"},
            {"S": "name"},
            {"S": "service"},
            {"S": "date"}
        ]},
        "parameterCount": {"N": "4"},
        "description": {"S": "Header: {{1}}=empresa | Body: 👋Hola {{2}}=nombre, recordatorio: tu {{3}}=servicio está programado para {{4}}=fecha según nuestro registro."},
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
        "externalId": {"S": "HX879326777b5252e74303734ca3b8066d"},
        "parameters": {"L": [
            {"S": "company"},
            {"S": "name"},
            {"S": "service"},
            {"S": "status"}
        ]},
        "parameterCount": {"N": "4"},
        "description": {"S": "Header: {{1}}=empresa | Body: 👋 ¡Hola {{2}}=nombre!, tenemos una actualización sobre tu {{3}}=servicio. 📍 Estado actual: {{4}}=estado. Te mantendremos informado(a)."},
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
        "externalId": {"S": "HXa771b1c388451770e47d0291001c54b5"},
        "parameters": {"L": [
            {"S": "company"},
            {"S": "name"},
            {"S": "service"},
            {"S": "detail"}
        ]},
        "parameterCount": {"N": "4"},
        "description": {"S": "Header: {{1}}=empresa | Body: 👋 Hola {{2}}=nombre, tu {{3}}=servicio ha sido confirmado(a). 📌 Detalle: {{4}}=detalle. Si necesitas más información, contáctanos."},
        "active": {"BOOL": true},
        "createdAt": {"S": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}
    }' > /dev/null

echo -e "${GREEN}✓${NC}"

echo ""
echo -e "${GREEN}✓ Templates creados exitosamente${NC}"
echo ""
echo "Templates disponibles:"
echo "  1. recordatorio_general (4 parámetros)"
echo "     - company: {{1}} del header (nombre empresa)"
echo "     - name: {{2}} del body (nombre destinatario)"
echo "     - service: {{3}} del body"
echo "     - date: {{4}} del body"
echo ""
echo "  2. actualizacion_estado (4 parámetros)"
echo "     - company: {{1}} del header (nombre empresa)"
echo "     - name: {{2}} del body (nombre destinatario)"
echo "     - service: {{3}} del body"
echo "     - status: {{4}} del body"
echo ""
echo "  3. confirmacion_general (4 parámetros)"
echo "     - company: {{1}} del header (nombre empresa)"
echo "     - name: {{2}} del body (nombre destinatario)"
echo "     - service: {{3}} del body"
echo "     - detail: {{4}} del body"
echo ""
echo "Ejemplo de uso:"
echo 'curl -X POST http://localhost:3000/v1/notifications/whatsapp \'
echo '  -H "X-API-Key: nfy_..." \'
echo '  -H "Content-Type: application/json" \'
echo '  -d '"'"'{
    "to": "+573001234567",
    "template_id": "recordatorio_general",
    "parameters": {
      "company": "Mi Empresa SAS",
      "name": "Yeison Peñaranda",
      "service": "Suscripcion Mensual",
      "date": "1 de Diciembre a las 10:00 AM"
    }
  }'"'"
