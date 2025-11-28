#!/bin/bash

echo "🚀 Iniciando DynamoDB Local..."

# Levantar contenedor
docker run -d \
  -p 8000:8000 \
  --name dynamodb-local \
  amazon/dynamodb-local

echo "⏳ Esperando que Dynamo arranque..."
sleep 3

echo "🗂️ Creando tabla NotificationService..."

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

echo "✔️ Tabla creada correctamente."
echo "📌 Listado de tablas:"
aws dynamodb list-tables --endpoint-url http://localhost:8000

echo "🎉 DynamoDB local está listo."
