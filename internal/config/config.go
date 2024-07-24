package config

import (
	"os"
)

type Config struct {
	PostgresConnStr   string
	DynamoDBTableName string
	DynamoDBRegion    string
}

func LoadConfig() (*Config, error) {

	return &Config{
		PostgresConnStr:   os.Getenv("POSTGRES_CONN_STR"),
		DynamoDBTableName: os.Getenv("DYNAMODB_TABLE_NAME"),
		DynamoDBRegion:    os.Getenv("AWS_REGION"),
	}, nil
}
