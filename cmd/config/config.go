package config

import (
	Fmt "fmt"
	OS "os"
	Strconv "strconv"
	Strings "strings"
	Sync "sync"
	Time "time"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
)

const (
	dbHost            = "DB_HOST"
	dbPort            = "DB_PORT"
	dbUsername        = "DB_USERNAME"
	dbPassword        = "DB_PASSWORD"
	dbDatabase        = "DB_DATABASE"
	dbSchema          = "DB_SCHEMA"
	jwtSecret         = "JWT_SECRET"
	googleClientID    = "GOOGLE_CLIENT_ID"
	brokerProtocol    = "BROKER_PROTOCOL"
	brokerAddress     = "BROKER_ADDRESS"
	brokerPort        = "BROKER_PORT"
	brokerUsername    = "BROKER_USERNAME"
	brokerPassword    = "BROKER_PASSWORD"
	applicationPort   = "APP_PORT"
	rateLimitRequests = "API_RATE_LIMIT_REQUESTS"
	rateLimitWindow   = "API_RATE_LIMIT_WINDOW"
	maxPageSize       = "API_MAX_PAGE_SIZE"
	allowedOrigins    = "CORS_ALLOWED_ORIGINS"
)

var (
	messageBrokerConfigurationInstance MessageBrokerConfiguration
	databaseConfigurationInstance      DatabaseConfiguration
	serverConfigurationInstance        ServerConfiguration
	messageBrokerConfigurationOnce     Sync.Once
	databaseConfigurationOnce          Sync.Once
	serverConfigurationOnce            Sync.Once
)

type ServerConfiguration struct {
	Addr              string
	JwtSecret         string
	GoogleClientID    string
	RateLimitRequests int
	RateLimitWindow   Time.Duration
	MaxPageSize       int
	AllowedOrigins    []string
}

type DatabaseConfiguration struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Schema   string
}

type MessageBrokerConfiguration struct {
	Protocol string
	Username string
	Password string
	Addr     string
	Port     string
}

func (database *DatabaseConfiguration) ToDNS() string {
	return Fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", database.Host, database.Port, database.Username, database.Password, database.Database, database.Schema)
}

func Database() *DatabaseConfiguration {
	databaseConfigurationOnce.Do(func() {
		databaseConfigurationInstance = DatabaseConfiguration{
			Host:     OS.Getenv(dbHost),
			Port:     OS.Getenv(dbPort),
			Username: OS.Getenv(dbUsername),
			Password: OS.Getenv(dbPassword),
			Database: OS.Getenv(dbDatabase),
			Schema:   OS.Getenv(dbSchema),
		}
	})
	return &databaseConfigurationInstance
}

func MessageBroker() *MessageBrokerConfiguration {
	messageBrokerConfigurationOnce.Do(func() {
		messageBrokerConfigurationInstance = MessageBrokerConfiguration{
			Addr:     OS.Getenv(brokerAddress),
			Port:     OS.Getenv(brokerPort),
			Username: OS.Getenv(brokerUsername),
			Password: OS.Getenv(brokerPassword),
			Protocol: OS.Getenv(brokerProtocol),
		}
	})
	return &messageBrokerConfigurationInstance
}

func NewServerConfiguration() *ServerConfiguration {
	serverConfigurationOnce.Do(func() {
		secret := OS.Getenv(jwtSecret)
		if secret == "" {
			// Use your provider if no secret is injected from infrastructure configuration
			provider := ApplicationMiddleware.NewSecretKeyProvider()
			generated, err := provider.GenerateKey()
			if err != nil {
				Logger.Log.Fatalf("Failed to generate fallback secret: %v", err)
			}
			secret = generated
		}
		serverConfigurationInstance = ServerConfiguration{
			Addr:              OS.Getenv(applicationPort),
			JwtSecret:         secret,
			GoogleClientID:    OS.Getenv(googleClientID),
			RateLimitRequests: envInt(rateLimitRequests, 60),
			RateLimitWindow:   envDuration(rateLimitWindow, Time.Minute),
			MaxPageSize:       envInt(maxPageSize, 50),
			AllowedOrigins:    envList(allowedOrigins),
		}
	})

	return &serverConfigurationInstance
}

func envInt(key string, fallback int) int {
	value, err := Strconv.Atoi(OS.Getenv(key))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func envDuration(key string, fallback Time.Duration) Time.Duration {
	value, err := Time.ParseDuration(OS.Getenv(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func envList(key string) []string {
	values := Strings.Split(OS.Getenv(key), ",")
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = Strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}
