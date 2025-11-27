package internal

import (
	"os"
)

type EnvirVariables struct {
	Env            string
	DynamoEndpoint string
	DynamoDBTable  string
}

func Environments() EnvirVariables {
	return EnvirVariables{
		DynamoEndpoint: os.Getenv("DYNAMO_ENDPOINT"),
		DynamoDBTable:  os.Getenv("TABLE_NAME"),
		Env:            os.Getenv("ENV"),
	}
}
