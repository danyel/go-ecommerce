package main

import (
	OS "os"

	Factory "github.com/danyel/ecommerce/cmd/factory/context"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Router "github.com/danyel/ecommerce/cmd/router"
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
	startApplicationContextFactory := Factory.InitializeStartApplicationContextFactory().StartMessageBroker()
	Router.NewAPIRouter(startApplicationContextFactory.ServerConfiguration(), startApplicationContextFactory.WebHandlerContextFactory()).Start()
}
