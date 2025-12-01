package main

import (
	"log"

	"github.com/danyel/ecommerce/cmd/broker"
	"github.com/danyel/ecommerce/cmd/config"
	"github.com/danyel/ecommerce/cmd/database"
	"github.com/danyel/ecommerce/cmd/router"
	"github.com/danyel/ecommerce/internal/category"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// project setup is done here ..
func main() {
	var err error
	var db *gorm.DB
	err = godotenv.Load()
	sc := config.NewServerConfiguration()
	dc := config.NewDatabaseConfiguration()
	bc := config.NewBrokerConfiguration()
	db, err = database.Connect(&dc)
	b := broker.NewBroker()
	if b.CreateConnection(&bc) != nil {
		log.Fatal(err)
	}
	b.RegisterConsumer(category.CategoryCreated, category.HandleCategoryCreated)
	b.RegisterConsumer(broker.QueueConfig{
		Topic: category.ExchangeCategory,
		Queue: "categories.category_created5",
	}, category.HandleCategoryCreated2)
	if err = b.Start(); err != nil {
		log.Println(err.Error())
	}
	r := router.ApiDefinition{
		SC:     &sc,
		DB:     db,
		Broker: b,
	}
	if err != nil {
		log.Fatal(err)
	}
	r.Run(r.ConfigRouter())
}
