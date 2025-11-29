package main

import (
	"log"

	"github.com/danyel/ecommerce/cmd/config"
	"github.com/danyel/ecommerce/cmd/database"
	"github.com/danyel/ecommerce/cmd/router"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// project setup is done here ..
func main() {
	var err error
	var connect *gorm.DB
	err = godotenv.Load()
	serverConfiguration := config.NewServerConfiguration()
	databaseConfiguration := config.NewDatabaseConfiguration()
	applicationConfiguration := config.NewApplicationConfiguration(serverConfiguration, databaseConfiguration)
	connect, err = database.Connect(&applicationConfiguration.DatabaseConfiguration)
	r := router.ApiDefinition{
		ServerConfiguration: &applicationConfiguration.ServerConfiguration,
		DB:                  connect,
	}
	if err != nil {
		log.Fatal(err)
	}
	r.Run(r.ConfigRouter())
}
