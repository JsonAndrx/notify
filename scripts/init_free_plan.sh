#!/bin/bash

# Script para inicializar el plan FREE en DynamoDB

# Configuración
TABLE_NAME="NotificationService"
ENDPOINT="http://localhost:8000"

echo "Creando plan FREE en DynamoDB..."

aws dynamodb put-item \
    --table-name $TABLE_NAME \
    --endpoint-url $ENDPOINT \
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
    }' \
    --return-consumed-capacity TOTAL

if [ $? -eq 0 ]; then
    echo "✓ Plan FREE creado exitosamente"
else
    echo "✗ Error al crear el plan FREE"
    exit 1
fi
