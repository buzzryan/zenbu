package config

import (
	"os"
)

type Config struct {
	JWSSigningKey string
	DynamoConfig
}

type DynamoConfig struct {
	// Endpoint is the endpoint for local DynamoDB. It is host:port.
	// In production, it should be empty.
	Endpoint  string
	TableName string
}

// LoadConfigFromEnv initializes the configuration from environment variables.
// TODO: use library such as godotenv to load configuration from .env file.
func LoadConfigFromEnv() Config {
	return Config{
		JWSSigningKey: os.Getenv("JWS_SIGNING_KEY"),
		DynamoConfig: DynamoConfig{
			Endpoint:  os.Getenv("DYNAMO_ENDPOINT"),
			TableName: os.Getenv("DYNAMO_TABLE_NAME"),
		},
	}
}
