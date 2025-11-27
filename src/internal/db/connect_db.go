package db

import (
	"context"
	"fmt"

	"notify-backend/internal"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"net/http"
)

func NewDynamoClient() (*dynamodb.Client, error) {
	ctx := context.TODO()
	env := internal.Environments()

	fmt.Println("======= Dynamo DEBUG START =======")

	// 1. Mostrar todo el ENV
	fmt.Printf("ENV.DynamoEndpoint = '%s'\n", env.DynamoEndpoint)

	// 2. Cargar config de AWS
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		fmt.Println("ERROR loading AWS config:", err)
		return nil, err
	}

	// 3. Mostrar región efectiva
	fmt.Println("AWS Region loaded:", cfg.Region)

	endpoint := env.DynamoEndpoint
	fmt.Printf("Using DynamoDB endpoint (raw): %s\n", endpoint)

	// 4. Validar si viene vacío
	if endpoint == "" {
		fmt.Println("Dynamo ENDPOINT is EMPTY → using AWS cloud Dynamo")
		fmt.Println("======= Dynamo DEBUG END =======")
		return dynamodb.NewFromConfig(cfg), nil
	}

	// 5. Testear conexión desde dentro del contenedor
	resp, err := http.Get(endpoint)
	fmt.Println("HTTP GET test to endpoint:", endpoint)
	if err != nil {
		fmt.Println("HTTP connection ERROR:", err)
	} else {
		fmt.Println("HTTP connection SUCCESS, status:", resp.StatusCode)
	}

	// 6. Crear cliente con endpoint override
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	fmt.Println("DynamoDB client created with BaseEndpoint override")
	fmt.Println("======= Dynamo DEBUG END =======")

	return client, nil
}

