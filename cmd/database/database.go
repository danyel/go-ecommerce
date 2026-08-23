package database

import (
	Fmt "fmt"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Postgres "gorm.io/driver/postgres"
	Database "gorm.io/gorm"
)

func Connect(databaseConfiguration *Configuration.DatabaseConfiguration) (*Database.DB, error) {
	connectionString := databaseConfiguration.ToDNS()
	connection, err := Database.Open(Postgres.Open(connectionString), &Database.Config{})
	Logger.Log.Info("Database connection established with DSN: %s", connectionString)
	if err != nil {
		Logger.Log.Fatal(Fmt.Sprintf("Failed to connect to the database: [%s]", connectionString), err)
		return nil, err
	}

	Logger.Log.Info("Successfully connected to the database[%s]\n", connectionString)
	return connection, nil
}
