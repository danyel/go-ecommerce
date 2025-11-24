package database

import (
	"fmt"
	"log"

	"github.com/dnoulet/ecommerce/cmd/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(database *config.DatabaseConfiguration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", database.Host, database.Port, database.Username, database.Password, database.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
