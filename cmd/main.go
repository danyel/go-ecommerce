package main

import (
	"log"
	"os"

	"github.com/danyel/ecommerce/cmd/broker"
	"github.com/danyel/ecommerce/cmd/config"
	"github.com/danyel/ecommerce/cmd/database"
	"github.com/danyel/ecommerce/cmd/router"
	shopping_basket "github.com/danyel/ecommerce/internal/shopping-basket"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// project setup is done here ..
func main() {
	var err error
	var db *gorm.DB
	var locations []string
	if os.Getenv("ENV") == "dev" {
		locations = []string{".env.dev"}
	} else {
		locations = []string{".env"}
	}
	err = godotenv.Load(locations...)
	sc := config.NewServerConfiguration()
	dc := config.NewDatabaseConfiguration()
	bc := config.NewBrokerConfiguration()
	db, err = database.Connect(&dc)
	b := broker.NewBroker()
	if b.CreateConnection(&bc) != nil {
		log.Fatal(err)
	}
	b.RegisterConsumer(shopping_basket.ShoppingBasketCreated, shopping_basket.HandleShoppingBasketCreated)
	b.RegisterConsumer(shopping_basket.ShoppingBasketUpdated, shopping_basket.HandleShoppingBasketUpdated)
	if err = b.Start(); err != nil {
		log.Println(err.Error())
	}
	r := router.ApiDefinition{
		SC:             &sc,
		DB:             db,
		EventPublisher: b,
	}
	if err != nil {
		log.Fatal(err)
	}
	r.Run(r.ConfigRouter())
}
