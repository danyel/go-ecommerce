package database

import (
	Fmt "fmt"
	Log "log"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Postgres "gorm.io/driver/postgres"
	Database "gorm.io/gorm"
)

func Connect(databaseConfiguration *Configuration.DatabaseConfiguration) (*Database.DB, error) {
	connectionString := databaseConfiguration.ToDNS()
	connection, err := Database.Open(Postgres.Open(connectionString), &Database.Config{})
	Log.Printf("Database connection established with DSN: %s", connectionString)
	if err != nil {
		Log.Fatal(Fmt.Sprintf("Failed to connect to the database: [%s]", connectionString), err)
		return nil, err
	}

	Log.Printf("Successfully connected to the database[%s]\n", connectionString)
	return connection, nil
}
