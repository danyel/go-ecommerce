package broker

import (
	Context "context"
	JSON "encoding/json"
	Fmt "fmt"

	Logger "github.com/danyel/ecommerce/cmd/logger"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	AMQP "github.com/rabbitmq/amqp091-go"
)

type QueueConfig struct {
	Topic      string
	Queue      string
	RoutingKey string
}

type MessageBroker struct {
	queueRegistries []QueueRegistry
	channel         *AMQP.Channel
	queueConfigs    []QueueConfig
	connection      *AMQP.Connection
}

func (messageBroker *MessageBroker) CreateConnection(messageBrokerConfiguration *Configuration.MessageBrokerConfiguration) error {
	brokerURL := Fmt.Sprintf("%s://%s:%s@%s:%s", messageBrokerConfiguration.Protocol, messageBrokerConfiguration.Username, messageBrokerConfiguration.Password, messageBrokerConfiguration.Addr, messageBrokerConfiguration.Port)
	Logger.Log.Info("connecting to %s", brokerURL)
	messageBrokerConnection, err := AMQP.Dial(brokerURL)
	if err != nil {
		return err
	}

	if messageBroker.channel, err = messageBrokerConnection.Channel(); err != nil {
		panic(err)
	}

	messageBroker.connection = messageBrokerConnection

	return nil
}

type HandlerFunc func([]byte) error

type QueueRegistry struct {
	HandlerFunc HandlerFunc
	QueueConfig QueueConfig
}

func (messageBroker *MessageBroker) RegisterConsumer(queueConfig QueueConfig, handlerFunc HandlerFunc) {
	if messageBroker.queueRegistries == nil {
		messageBroker.queueRegistries = []QueueRegistry{{handlerFunc, queueConfig}}
	}
	notFound := true
	for _, registry := range messageBroker.queueRegistries {
		if registry.QueueConfig.Queue == queueConfig.Queue && registry.QueueConfig.Topic == queueConfig.Topic && registry.QueueConfig.RoutingKey == queueConfig.RoutingKey {
			notFound = false
			break
		}
	}

	if notFound {
		messageBroker.queueRegistries = append(messageBroker.queueRegistries, QueueRegistry{handlerFunc, queueConfig})
	}
}

func (messageBroker *MessageBroker) setup() error {
	for _, c := range messageBroker.queueRegistries {
		if !messageBroker.alreadyRegistered(c.QueueConfig) {
			err := messageBroker.registerQueue(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (messageBroker *MessageBroker) registerQueue(queueRegistry QueueRegistry) error {
	if err := messageBroker.channel.ExchangeDeclare(queueRegistry.QueueConfig.Topic, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := messageBroker.channel.QueueDeclare(queueRegistry.QueueConfig.Queue, true, false, false, false, nil); err != nil {
		return err
	}
	if err := messageBroker.channel.QueueBind(queueRegistry.QueueConfig.Queue, queueRegistry.QueueConfig.RoutingKey, queueRegistry.QueueConfig.Topic, false, nil); err != nil {
		return err
	}
	messageBroker.queueConfigs = append(messageBroker.queueConfigs, queueRegistry.QueueConfig)
	return nil
}

func (messageBroker *MessageBroker) alreadyRegistered(queueConfig QueueConfig) bool {
	found := false
	for _, currentQueueConfig := range messageBroker.queueConfigs {
		if queueConfig.Topic == currentQueueConfig.Topic && queueConfig.Queue == currentQueueConfig.Queue && queueConfig.RoutingKey == currentQueueConfig.RoutingKey {
			found = true
			break
		}
	}

	return found
}

func (messageBroker *MessageBroker) Publish(queue string, value any) error {
	Logger.Log.Debug("Publishing message to queue %s with value %v", queue, value)
	for _, queueRegistry := range messageBroker.queueRegistries {
		Logger.Log.Info("Registered: %v", messageBroker)
		if queueRegistry.QueueConfig.Queue == queue && messageBroker.alreadyRegistered(queueRegistry.QueueConfig) {
			body, e := JSON.Marshal(value)
			if e != nil {
				return e
			}
			return messageBroker.channel.PublishWithContext(Context.Background(), queueRegistry.QueueConfig.Topic, queueRegistry.QueueConfig.RoutingKey, false, false, AMQP.Publishing{ContentType: "application/json", Body: body})
		}
	}
	return Fmt.Errorf("no handler registered for queue %s", queue)
}

func (messageBroker *MessageBroker) consume(queueRegistry QueueRegistry) {
	var err error
	var messages <-chan AMQP.Delivery
	if messages, err = messageBroker.channel.Consume(queueRegistry.QueueConfig.Queue, "", false, false, false, false, nil); err != nil {
		Logger.Log.Fatalf("Error on consuming message: %v", err.Error())
	}
	go func() {
		for message := range messages {
			if err := queueRegistry.HandlerFunc(message.Body); err != nil {
				_ = message.Nack(false, false)
				continue
			}
			_ = message.Ack(false)
		}
	}()
	Logger.Log.Debug("[Consumer] Listening: %v", queueRegistry.QueueConfig.Queue)
}

func (messageBroker *MessageBroker) Start() error {
	if err := messageBroker.setup(); err != nil {
		return err
	}
	for _, queueRegistry := range messageBroker.queueRegistries {
		if messageBroker.alreadyRegistered(queueRegistry.QueueConfig) {
			go messageBroker.consume(queueRegistry)
		}
	}

	context, cancel := Context.WithCancel(Context.Background())
	go func() {
		for {
			select {
			case <-context.Done():
				return
			default:
				for _, queueRegistry := range messageBroker.queueRegistries {
					if !messageBroker.alreadyRegistered(queueRegistry.QueueConfig) {
						err := messageBroker.registerQueue(queueRegistry)

						if err != nil {
							Logger.Log.Fatalf("Error on registering registry: %s", err.Error())
						} else {
							messageBroker.queueConfigs = append(messageBroker.queueConfigs, queueRegistry.QueueConfig)
						}
					}
				}
				if len(messageBroker.queueConfigs) == len(messageBroker.queueRegistries) {
					cancel()
				}
			}
		}
	}()

	return nil
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		queueRegistries: make([]QueueRegistry, 0),
		queueConfigs:    make([]QueueConfig, 0),
	}
}
