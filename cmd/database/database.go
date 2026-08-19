package database

import (
	Fmt "fmt"
	Log "log"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Postgres "gorm.io/driver/postgres"
	Database "gorm.io/gorm"
)

func Connect(database *Configuration.DatabaseConfiguration) (*Database.DB, error) {
	dns := database.ToDNS()
	db, err := Database.Open(Postgres.Open(dns), &Database.Config{})
	Log.Printf("Database connection established with DSN: %s", dns)
	if err != nil {
		Log.Fatal(Fmt.Sprintf("Failed to connect to the database: [%s]", dns), err)
		return nil, err
	}

	Log.Printf("Successfully connected to the database[%s]\n", dns)
	return db, nil
}
