package factory

import (
	Errors "errors"

	MessageBroker "github.com/danyel/ecommerce/cmd/broker"
	Configuration "github.com/danyel/ecommerce/cmd/config"
	Logger "github.com/danyel/ecommerce/cmd/logger"
)

type MessageBrokerConnectionFactory interface {
	MessageBroker() *MessageBroker.MessageBroker
	Start() error
}

type messageBrokerConnectionFactory struct {
	messageBrokerConfiguration *Configuration.MessageBrokerConfiguration
	messageBroker              *MessageBroker.MessageBroker
}

func (messageBrokerConnectionFactory *messageBrokerConnectionFactory) MessageBroker() *MessageBroker.MessageBroker {
	if messageBrokerConnectionFactory.messageBroker == nil {
		Logger.Log.Fatal("Critical: MessageBroker returned nil from the connection factory")
	}
	return messageBrokerConnectionFactory.messageBroker
}

func (messageBrokerConnectionFactory *messageBrokerConnectionFactory) Start() error {
	return messageBrokerConnectionFactory.messageBroker.Start()
}

func createMessageBroker(messageBrokerConfiguration *Configuration.MessageBrokerConfiguration) *MessageBroker.MessageBroker {
	Logger.Log.Debug("[BrokerConnectionFactory] Creating broker")
	messageBroker := MessageBroker.NewMessageBroker()

	if messageBroker.CreateConnection(messageBrokerConfiguration) != nil {
		Logger.Log.Fatal(Errors.New("Can not connect to broker: %s" + messageBrokerConfiguration.Addr))
	}

	return messageBroker
}

func NewMessageBrokerConnectionFactory(messageBrokerConfiguration *Configuration.MessageBrokerConfiguration) MessageBrokerConnectionFactory {
	brokerConfigurationOnce.Do(func() {
		brokerConnectionFactoryInstance = &messageBrokerConnectionFactory{
			messageBrokerConfiguration: messageBrokerConfiguration,
			messageBroker:              createMessageBroker(messageBrokerConfiguration),
		}
	})
	return brokerConnectionFactoryInstance
}
