package main

import (
	OS "os"

	Logger "github.com/danyel/ecommerce/cmd/logger"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Router "github.com/danyel/ecommerce/cmd/router"
	Factory "github.com/danyel/ecommerce/internal/common/factory"
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
	}
	serverConfiguration := Configuration.NewServerConfiguration()
	databaseConfiguration := Configuration.NewDatabaseConfiguration()
	messageBrokerConfiguration := Configuration.NewMessageBrokerConfiguration()
	databaseConnectionFactory := Factory.NewDatabaseConnectionFactory(&databaseConfiguration)
	brokerConnectionFactory := Factory.NewMessageBrokerConnectionFactory(&messageBrokerConfiguration)
	applicationConnectionFactory := Factory.NewApplicationConnectionFactory(databaseConnectionFactory, brokerConnectionFactory)

	ShoppingBasket.RegisterShoppingBasketEvents(applicationConnectionFactory.ShoppingBasketService(), brokerConnectionFactory.MessageBroker())
	Reservation.RegisterReservationEvents(applicationConnectionFactory.ReservationService(), applicationConnectionFactory.ProductService(), brokerConnectionFactory.MessageBroker())
	if err = brokerConnectionFactory.MessageBroker().Start(); err != nil {
		Logger.Log.Debug("%v", err.Error())
	}
	apiDefinition := Router.APIDefinition{
		ServerConfiguration:          &serverConfiguration,
		ApplicationConnectionFactory: applicationConnectionFactory,
	}
	if err != nil {
		Logger.Log.Fatal(err)
	}
	apiDefinition.Run(apiDefinition.ConfigRouter())
}
