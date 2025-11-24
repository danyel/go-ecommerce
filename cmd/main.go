package main

import (
	"log"

	"github.com/dnoulet/ecommerce/cmd/config"
	"github.com/dnoulet/ecommerce/cmd/database"
	"github.com/dnoulet/ecommerce/cmd/router"
)

// project setup is done here ..
func main() {
	serverConfiguration := config.NewServerConfiguration()
	databaseConfiguration := config.NewDatabaseConfiguration()
	applicationConfiguration := config.NewApplicationConfiguration(serverConfiguration, databaseConfiguration)
	connect, err := database.Connect(&applicationConfiguration.DatabaseConfiguration)
	r := router.ApiDefinition{
		ServerConfiguration: &applicationConfiguration.ServerConfiguration,
		DB:                  connect,
	}
	if err != nil {
		log.Fatal(err)
	}
	r.Run(r.ConfigRouter())
}
