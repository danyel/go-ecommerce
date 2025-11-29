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

type ApplicationConfiguration struct {
	DatabaseConfiguration DatabaseConfiguration
	ServerConfiguration   ServerConfiguration
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

func NewApplicationConfiguration(s ServerConfiguration, d DatabaseConfiguration) ApplicationConfiguration {
	return ApplicationConfiguration{
		ServerConfiguration:   s,
		DatabaseConfiguration: d,
	}
}
