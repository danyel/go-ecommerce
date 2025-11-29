package database

import (
	"fmt"
	"log"

	"github.com/danyel/ecommerce/cmd/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(database *config.DatabaseConfiguration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", database.Host, database.Port, database.Username, database.Password, database.Database, database.Schema)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Printf("Database connection established with DSN: %s", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
		return nil, err
	}

	log.Printf("Successfully connected to the database[%s]\n", dsn)
	return db, nil
}
