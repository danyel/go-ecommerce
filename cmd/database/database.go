package database

import (
	Fmt "fmt"
	Log "log"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Postgres "gorm.io/driver/postgres"
	Database "gorm.io/gorm"
)

func Connect(database *Configuration.DatabaseConfiguration) (*Database.DB, error) {
	dsn := Fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", database.Host, database.Port, database.Username, database.Password, database.Database, database.Schema)
	db, err := Database.Open(Postgres.Open(dsn), &Database.Config{})
	Log.Printf("Database connection established with DSN: %s", dsn)
	if err != nil {
		Log.Fatal(Fmt.Sprintf("Failed to connect to the database: [%s]", dsn), err)
		return nil, err
	}

	Log.Printf("Successfully connected to the database[%s]\n", dsn)
	return db, nil
}
