package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWSSigningKey string
	MySQLConfig
	DynamoConfig
}

type MySQLConfig struct {
	// Endpoint is `hostname:port`
	Endpoint   string
	User       string
	Password   string
	DBName     string
	LogEnabled bool
}

func (m *MySQLConfig) DSN() string {
	return m.User + ":" + m.Password + "@tcp(" + m.Endpoint + ")/" + m.DBName + "?parseTime=true&loc=UTC&charset=utf8mb4"
}

type DynamoConfig struct {
	// Endpoint is the endpoint for local DynamoDB. It is host:port.
	// In production, it should be empty.
	Endpoint string
}

// LoadConfigFromEnv initializes the configuration from environment variables.
// TODO: use library such as godotenv to load configuration from .env file.
func LoadConfigFromEnv() Config {
	mysqlLogEnabledString := os.Getenv("MYSQL_LOG_ENABLED")
	// If MYSQL_LOG_ENABLED is not set or not valid boolean, it is considered as false.
	mysqlLogEnabled, _ := strconv.ParseBool(mysqlLogEnabledString)
	return Config{
		JWSSigningKey: os.Getenv("JWS_SIGNING_KEY"),
		MySQLConfig: MySQLConfig{
			Endpoint:   os.Getenv("MYSQL_ENDPOINT"),
			User:       os.Getenv("MYSQL_USER"),
			Password:   os.Getenv("MYSQL_PASSWORD"),
			DBName:     os.Getenv("MYSQL_DATABASE"),
			LogEnabled: mysqlLogEnabled,
		},
		DynamoConfig: DynamoConfig{
			Endpoint: os.Getenv("DYNAMO_ENDPOINT"),
		},
	}
}
