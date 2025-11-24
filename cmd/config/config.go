package config

type ServerConfiguration struct {
	Addr string
}

type DatabaseConfiguration struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type ApplicationConfiguration struct {
	DatabaseConfiguration DatabaseConfiguration
	ServerConfiguration   ServerConfiguration
}

// TODO this properties needs to be read from .env
func NewDatabaseConfiguration() DatabaseConfiguration {
	return DatabaseConfiguration{
		Host:     "localhost",
		Port:     "5401",
		Username: "ecommerce",
		Password: "ecommerce",
		Database: "ecommerce",
	}
}

// TODO this properties needs to be read from .env
func NewServerConfiguration() ServerConfiguration {
	return ServerConfiguration{
		Addr: ":8080",
	}
}

func NewApplicationConfiguration(s ServerConfiguration, d DatabaseConfiguration) ApplicationConfiguration {
	return ApplicationConfiguration{
		ServerConfiguration:   s,
		DatabaseConfiguration: d,
	}
}
