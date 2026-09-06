package contextfactory

import (
	OS "os"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Reservation "github.com/danyel/ecommerce/internal/product"
)

type StartApplicationContextFactory interface {
	StartMessageBroker() StartApplicationContextFactory
	WebHandlerContextFactory() WebHandlerContextFactory
	ServerConfiguration() *Configuration.ServerConfiguration
}

type startApplicationContextFactory struct {
	startMessageBrokerFunc   func() error
	webHandlerContextFactory WebHandlerContextFactory
	serviceContextFactory    ServiceContextFactory
}

func (startApplicationContextFactory *startApplicationContextFactory) registerEvents() {
	Reservation.RegisterConsumer(startApplicationContextFactory.serviceContextFactory.ProductService(), startApplicationContextFactory.serviceContextFactory.MessageBroker())
}

func (startApplicationContextFactory *startApplicationContextFactory) StartMessageBroker() StartApplicationContextFactory {
	if err := startApplicationContextFactory.startMessageBrokerFunc(); err != nil {
		Logger.Log.Debug("%v", err.Error())
		OS.Exit(0)
	}
	return startApplicationContextFactory
}

func (startApplicationContextFactory *startApplicationContextFactory) WebHandlerContextFactory() WebHandlerContextFactory {
	return startApplicationContextFactory.webHandlerContextFactory
}

func (startApplicationContextFactory *startApplicationContextFactory) ServerConfiguration() *Configuration.ServerConfiguration {
	return Configuration.NewServerConfiguration()
}

func InitializeStartApplicationContextFactory() StartApplicationContextFactory {
	Logger.Log.Info("Application Initialization started")
	InitializeDatabaseContextFactory()
	InitializeMessageBrokerContextFactory()
	BuildServiceContextFactory()
	webHandlerContextFactory := BuildAWebHandlerContextFactory()
	var startApplicationContextFactoryInstance = &startApplicationContextFactory{
		webHandlerContextFactory: webHandlerContextFactory,
		serviceContextFactory:    applicationContextFactoryInstance,
		startMessageBrokerFunc: func() error {
			err := applicationContextFactoryInstance.StartMessageBroker()
			return err
		},
	}
	startApplicationContextFactoryInstance.registerEvents()
	Logger.Log.Info("Application Initialization finished")
	return startApplicationContextFactoryInstance
}
