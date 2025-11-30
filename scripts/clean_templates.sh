#!/bin/bash

# Script para eliminar templates existentes de DynamoDB

set -e

TABLE_NAME="NotificationService"
ENDPOINT="http://localhost:8000"

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "🗑️  Eliminando Templates Existentes"
echo "===================================="
echo ""

# Template 1: recordatorio_general
echo -n "Eliminando template: recordatorio_general... "
aws dynamodb delete-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
    --region us-east-1 \
    --key '{
        "PK": {"S": "TEMPLATE#recordatorio_general"},
        "SK": {"S": "METADATA"}
    }' > /dev/null 2>&1 || true

echo -e "${GREEN}✓${NC}"

# Template 2: actualizacion_estado
echo -n "Eliminando template: actualizacion_estado... "
aws dynamodb delete-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
    --region us-east-1 \
    --key '{
        "PK": {"S": "TEMPLATE#actualizacion_estado"},
        "SK": {"S": "METADATA"}
    }' > /dev/null 2>&1 || true

echo -e "${GREEN}✓${NC}"

# Template 3: confirmacion_general
echo -n "Eliminando template: confirmacion_general... "
aws dynamodb delete-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
    --region us-east-1 \
    --key '{
        "PK": {"S": "TEMPLATE#confirmacion_general"},
        "SK": {"S": "METADATA"}
    }' > /dev/null 2>&1 || true

echo -e "${GREEN}✓${NC}"

echo ""
echo -e "${GREEN}✓ Templates eliminados exitosamente${NC}"
echo ""
echo "Ahora puedes ejecutar: ./init_templates.sh"
