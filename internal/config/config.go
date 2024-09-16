package config

import (
	"os"
)

type Config struct {
	JWSSigningKey string
	DynamoConfig
	S3Config
}

type DynamoConfig struct {
	// Endpoint is the endpoint for local DynamoDB. It is host:port.
	// In production, it should be empty.
	Endpoint  string
	TableName string
}

type S3Config struct {
	Bucket                   string
	PrivateDir               string
	PublicDir                string
	PublicCloudfrontEndpoint string
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
		S3Config: S3Config{
			Bucket:                   os.Getenv("S3_BUCKET"),
			PrivateDir:               os.Getenv("S3_PRIVATE_DIR"),
			PublicDir:                os.Getenv("S3_PUBLIC_DIR"),
			PublicCloudfrontEndpoint: os.Getenv("S3_PUBLIC_CLOUDFRONT_ENDPOINT"),
		},
	}
}
