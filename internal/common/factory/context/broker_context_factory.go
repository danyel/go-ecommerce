package contextfactory

import (
	Errors "errors"

	MessageBroker "github.com/danyel/ecommerce/cmd/broker"
	Configuration "github.com/danyel/ecommerce/cmd/config"
	Logger "github.com/danyel/ecommerce/cmd/logger"
)

type MessageBrokerContextFactory interface {
	Start() error
	MessageBroker() *MessageBroker.MessageBroker
}

type messageBrokerContextFactory struct {
	messageBrokerConfiguration *Configuration.MessageBrokerConfiguration
	messageBroker              *MessageBroker.MessageBroker
}

func (messageBrokerContextFactory *messageBrokerContextFactory) Start() error {
	return messageBrokerContextFactory.messageBroker.Start()
}

func createMessageBroker(messageBrokerConfiguration *Configuration.MessageBrokerConfiguration) *MessageBroker.MessageBroker {
	Logger.Log.Debug("[BrokerConnectionFactory] Creating broker")
	messageBroker := MessageBroker.NewMessageBroker()

	if messageBroker.CreateConnection(messageBrokerConfiguration) != nil {
		Logger.Log.Fatal(Errors.New("Can not connect to broker: %s" + messageBrokerConfiguration.Addr))
	}

	return messageBroker
}

func (messageBrokerContextFactory *messageBrokerContextFactory) MessageBroker() *MessageBroker.MessageBroker {
	if messageBrokerContextFactory.getMessageBroker() == nil {
		Logger.Log.Fatal("No message broker and message broker configuration found")
	}
	return messageBrokerContextFactory.messageBroker
}

func (messageBrokerContextFactory *messageBrokerContextFactory) getMessageBroker() *MessageBroker.MessageBroker {
	if messageBrokerContextFactory.messageBroker != nil {
		return messageBrokerContextFactory.messageBroker
	}
	if messageBrokerContextFactory.messageBrokerConfiguration != nil {
		return createMessageBroker(messageBrokerContextFactory.messageBrokerConfiguration)
	}
	return nil
}

func InitializeMessageBrokerContextFactory(messageBrokerConfiguration *Configuration.MessageBrokerConfiguration) {
	messageBrokerContextFactoryOnce.Do(func() {
		messageBrokerContextFactoryInstance = &messageBrokerContextFactory{
			messageBrokerConfiguration: messageBrokerConfiguration,
			messageBroker:              createMessageBroker(messageBrokerConfiguration),
		}
	})
}
