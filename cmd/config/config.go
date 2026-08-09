package config

import (
	Log "log"
	OS "os"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
)

type ServerConfiguration struct {
	Addr      string
	JwtSecret string
}

type DatabaseConfiguration struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Schema   string
}

type BrokerConfiguration struct {
	Protocol string
	Username string
	Password string
	Addr     string
	Port     string
}

func NewDatabaseConfiguration() DatabaseConfiguration {
	return DatabaseConfiguration{
		Host:     OS.Getenv("DB_HOST"),
		Port:     OS.Getenv("DB_PORT"),
		Username: OS.Getenv("DB_USERNAME"),
		Password: OS.Getenv("DB_PASSWORD"),
		Database: OS.Getenv("DB_DATABASE"),
		Schema:   OS.Getenv("DB_SCHEMA"),
	}
}

func NewServerConfiguration() ServerConfiguration {
	secret := OS.Getenv("JWT_SECRET")
	if secret == "" {
		// Use your provider if no secret is injected from infrastructure configuration
		provider := ApplicationMiddleware.NewSecretKeyProvider()
		generated, err := provider.GenerateKey()
		if err != nil {
			Log.Fatalf("Failed to generate fallback secret: %v", err)
		}
		secret = generated
	}
	return ServerConfiguration{
		Addr:      OS.Getenv("APP_PORT"),
		JwtSecret: secret,
	}
}

func NewBrokerConfiguration() BrokerConfiguration {
	return BrokerConfiguration{
		Protocol: OS.Getenv("BROKER_PROTOCOL"),
		Addr:     OS.Getenv("BROKER_ADDRESS"),
		Port:     OS.Getenv("BROKER_PORT"),
		Username: OS.Getenv("BROKER_USERNAME"),
		Password: OS.Getenv("BROKER_PASSWORD"),
	}
}
