package config

import "os"

type ServerConfiguration struct {
	Addr string
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
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Schema:   os.Getenv("DB_SCHEMA"),
	}
}

func NewServerConfiguration() ServerConfiguration {
	return ServerConfiguration{
		Addr: os.Getenv("APP_PORT"),
	}
}

func NewBrokerConfiguration() BrokerConfiguration {
	return BrokerConfiguration{
		Protocol: os.Getenv("BROKER_PROTOCOL"),
		Addr:     os.Getenv("BROKER_ADDRESS"),
		Port:     os.Getenv("BROKER_PORT"),
		Username: os.Getenv("BROKER_USERNAME"),
		Password: os.Getenv("BROKER_PASSWORD"),
	}
}
