package database

import (
	Fmt "fmt"
	Regex "regexp"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Postgres "gorm.io/driver/postgres"
	Database "gorm.io/gorm"
)

var passwordRegex = Regex.MustCompile(`(?i)(password=)([^ ]+)`)

func Connect(databaseConfiguration *Configuration.DatabaseConfiguration) (*Database.DB, error) {
	connectionString := databaseConfiguration.ToDNS()
	connection, err := Database.Open(Postgres.Open(connectionString), &Database.Config{})
	Logger.Log.Info("Database connection established with DSN: %s", markPasswordObfuscating(connectionString))
	if err != nil {
		Logger.Log.Fatal(Fmt.Sprintf("Failed to connect to the database: [%s]", markPasswordObfuscating(connectionString)), err)
		return nil, err
	}

	Logger.Log.Info("⚠ Successfully connected to the database!!🥳")
	return connection, nil
}

func markPasswordObfuscating(value string) string {
	return passwordRegex.ReplaceAllString(value, "${1}****")
}
