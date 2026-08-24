package main

import (
	OS "os"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Router "github.com/danyel/ecommerce/cmd/router"
	Factory "github.com/danyel/ecommerce/internal/common/factory/context"
	Reservation "github.com/danyel/ecommerce/internal/reservation"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shoppingbasket"
	GoDotEnv "github.com/joho/godotenv"
)

// project setup is done here.
func main() {
	var err error
	var locations []string
	if OS.Getenv("ENV") == "dev" {
		locations = []string{".env.dev"}
	} else {
		locations = []string{".env"}
	}
	err = GoDotEnv.Load(locations...)
	if err != nil {
		Logger.Log.Fatal(err)
		OS.Exit(0)
	}
	serverConfiguration := Configuration.NewServerConfiguration()
	databaseConfiguration := Configuration.NewDatabaseConfiguration()
	messageBrokerConfiguration := Configuration.NewMessageBrokerConfiguration()
	Factory.InitializeDatabaseContextFactory(&databaseConfiguration)
	Factory.InitializeMessageBrokerContextFactory(&messageBrokerConfiguration)
	applicationContextFactory := Factory.BuildApplicationContextFactory()
	webHandlerContextFactory := Factory.BuildAWebHandlerContextFactory()

	ShoppingBasket.RegisterShoppingBasketEvents(applicationContextFactory.ShoppingBasketService(), applicationContextFactory.MessageBroker())
	Reservation.RegisterReservationEvents(applicationContextFactory.ReservationService(), applicationContextFactory.ProductService(), applicationContextFactory.MessageBroker())
	if err = applicationContextFactory.StartMessageBroker(); err != nil {
		Logger.Log.Debug("%v", err.Error())
		OS.Exit(0)
	}
	Router.NewApiRouter(&serverConfiguration, webHandlerContextFactory).Start()
}
