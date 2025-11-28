#!/bin/bash

# Script de setup completo para el proyecto Notification Service

set -e

echo "🚀 Notification Service - Setup"
echo "================================"
echo ""

# Configuración
TABLE_NAME="NotificationService"
ENDPOINT="http://localhost:8000"

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Función para verificar si DynamoDB está corriendo
check_dynamodb() {
    echo -n "Verificando DynamoDB local... "
    if aws dynamodb list-tables --endpoint-url $ENDPOINT --region us-east-1 &>/dev/null; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        return 1
    fi
}

# Verificar DynamoDB
if ! check_dynamodb; then
    echo -e "${YELLOW}⚠️  DynamoDB local no está corriendo${NC}"
    echo "Por favor, ejecuta: docker run -p 8000:8000 amazon/dynamodb-local"
    exit 1
fi

# Verificar si la tabla existe
echo -n "Verificando tabla $TABLE_NAME... "
if aws dynamodb describe-table --table-name $TABLE_NAME --endpoint-url $ENDPOINT --region us-east-1 &>/dev/null; then
    echo -e "${YELLOW}Ya existe${NC}"
else
    echo -e "${RED}No existe${NC}"
    echo "Creando tabla..."
    
    aws dynamodb create-table \
        --table-name $TABLE_NAME \
        --attribute-definitions \
            AttributeName=PK,AttributeType=S \
            AttributeName=SK,AttributeType=S \
        --key-schema \
            AttributeName=PK,KeyType=HASH \
            AttributeName=SK,KeyType=RANGE \
        --billing-mode PAY_PER_REQUEST \
        --endpoint-url $ENDPOINT \
        --region us-east-1 > /dev/null
    
    echo -e "${GREEN}✓ Tabla creada${NC}"
fi

# Verificar si el plan FREE existe
echo -n "Verificando plan FREE... "
if aws dynamodb get-item \
    --table-name $TABLE_NAME \
    --key '{"PK": {"S": "PLAN#FREE"}, "SK": {"S": "METADATA"}}' \
    --endpoint-url $ENDPOINT \
    --region us-east-1 &>/dev/null; then
    echo -e "${YELLOW}Ya existe${NC}"
else
    echo -e "${RED}No existe${NC}"
    echo "Creando plan FREE..."
    
    aws dynamodb put-item \
        --table-name $TABLE_NAME \
        --endpoint-url $ENDPOINT \
        --region us-east-1 \
        --item '{
            "PK": {"S": "PLAN#FREE"},
            "SK": {"S": "METADATA"},
            "name": {"S": "Free Plan"},
            "notificationLimit": {"N": "50"},
            "periodDays": {"N": "30"},
            "price": {"N": "0.00"},
            "description": {"S": "Plan gratuito con 50 notificaciones cada 30 días"},
            "active": {"BOOL": true},
            "createdAt": {"S": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}
        }' > /dev/null
    
    echo -e "${GREEN}✓ Plan FREE creado${NC}"
fi

echo ""
echo -e "${GREEN}✓ Setup completado exitosamente${NC}"
echo ""
echo "Próximos pasos:"
echo "1. cd infrastructure"
echo "2. sam build"
echo "3. sam local start-api --env-vars env.json"
echo ""
echo "La API estará disponible en: http://localhost:3000"
