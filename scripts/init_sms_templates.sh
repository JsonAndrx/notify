#!/bin/bash

# Script para inicializar templates de SMS en DynamoDB

set -e

TABLE_NAME="NotificationService"
ENDPOINT="http://localhost:8000"

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "🔧 Inicializando Templates de SMS"
echo "==================================="
echo ""

# Template 1: Código de Verificación (SMS)
echo -n "Creando template: sms_verification_code... "
aws dynamodb put-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
    --region us-east-1 \
    --item '{
        "PK": {"S": "TEMPLATE#sms_verification_code"},
        "SK": {"S": "METADATA"},
        "templateId": {"S": "sms_verification_code"},
        "name": {"S": "Código de Verificación SMS"},
        "type": {"S": "sms"},
        "provider": {"S": "twilio"},
        "externalId": {"S": ""},
        "parameters": {"L": [
            {"S": "codigo"}
        ]},
        "parameterCount": {"N": "1"},
        "description": {"S": "Su codigo de verificacion para $empresa es: $codigo"},
        "active": {"BOOL": true},
        "createdAt": {"S": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}
    }' > /dev/null

echo -e "${GREEN}✓${NC}"

echo ""
echo -e "${GREEN}✓ Template de SMS creado exitosamente${NC}"
echo ""
echo "Template disponible:"
echo "  - sms_verification_code"
echo "    • codigo: Código numérico de 4-6 dígitos"
echo "    • empresa: Se toma automáticamente del negocio"
echo ""
echo "Formato del mensaje:"
echo "  Su codigo de verificacion para [EMPRESA] es: [CODIGO]"
echo ""
echo "Ejemplo de uso:"
echo 'curl -X POST http://localhost:3000/v1/notifications/sms \'
echo '  -H "X-API-Key: nfy_..." \'
echo '  -H "Content-Type: application/json" \'
echo '  -d '"'"'{
    "to": "+573001234567",
    "template_id": "sms_verification_code",
    "parameters": {
      "codigo": "123456"
    }
  }'"'"
