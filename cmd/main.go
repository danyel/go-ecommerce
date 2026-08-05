package main

import (
	Log "log"
	OS "os"

	Broker "github.com/danyel/ecommerce/cmd/broker"
	Config "github.com/danyel/ecommerce/cmd/config"
	Database "github.com/danyel/ecommerce/cmd/database"
	Router "github.com/danyel/ecommerce/cmd/router"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shopping-basket"
	GoDotEnv "github.com/joho/godotenv"
	Gorm "gorm.io/gorm"
)

// project setup is done here.
func main() {
	var err error
	var db *Gorm.DB
	var locations []string
	if OS.Getenv("ENV") == "dev" {
		locations = []string{".env.dev"}
	} else {
		locations = []string{".env"}
	}
	err = GoDotEnv.Load(locations...)
	sc := Config.NewServerConfiguration()
	dc := Config.NewDatabaseConfiguration()
	bc := Config.NewBrokerConfiguration()
	db, err = Database.Connect(&dc)
	b := Broker.NewBroker()
	if b.CreateConnection(&bc) != nil {
		Log.Fatal(err)
	}
	shoppingBasketEventHandler := ShoppingBasket.NewShoppingBasketEventHandler(ShoppingBasket.NewService(db, b))
	b.RegisterConsumer(ShoppingBasket.ShoppingBasketCreated, shoppingBasketEventHandler.HandleShoppingBasketCreated)
	b.RegisterConsumer(ShoppingBasket.ShoppingBasketUpdated, shoppingBasketEventHandler.HandleShoppingBasketUpdated)
	if err = b.Start(); err != nil {
		Log.Println(err.Error())
	}
	r := Router.ApiDefinition{
		SC:             &sc,
		DB:             db,
		EventPublisher: b,
	}
	if err != nil {
		Log.Fatal(err)
	}
	r.Run(r.ConfigRouter())
}
