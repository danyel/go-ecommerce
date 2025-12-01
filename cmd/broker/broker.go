package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/danyel/ecommerce/cmd/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueConfig struct {
	Topic      string
	Queue      string
	RoutingKey string
}

type Broker struct {
	registry   []QueueRegistry
	channel    *amqp.Channel
	registered []QueueConfig
}

func (b *Broker) CreateConnection(c *config.BrokerConfiguration) error {
	brokerUrl := fmt.Sprintf("%s://%s:%s@%s:%s", c.Protocol, c.Username, c.Password, c.Addr, c.Port)
	slog.Info("connecting to " + brokerUrl)
	connection, err := amqp.Dial(brokerUrl)
	if err != nil {
		return err
	}

	if b.channel, err = connection.Channel(); err != nil {
		panic(err)
	}

	return nil
}

type HandlerFunc func([]byte) error

type QueueRegistry struct {
	H HandlerFunc
	C QueueConfig
}

func (b *Broker) RegisterConsumer(queue QueueConfig, handler HandlerFunc) {
	if b.registry == nil {
		b.registry = []QueueRegistry{{handler, queue}}
	}
	notFound := true
	for _, registry := range b.registry {
		if registry.C.Queue == queue.Queue && registry.C.Topic == queue.Topic && registry.C.RoutingKey == queue.RoutingKey {
			notFound = false
			break
		}
	}

	if notFound {
		b.registry = append(b.registry, QueueRegistry{handler, queue})
	}
}

func (b *Broker) setup() error {
	for _, c := range b.registry {
		if !b.alreadyRegistered(c.C) {
			err := b.register(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Broker) register(c QueueRegistry) error {
	if er := b.channel.ExchangeDeclare(c.C.Topic, "topic", true, false, false, false, nil); er != nil {
		return er
	}
	if _, er := b.channel.QueueDeclare(c.C.Queue, true, false, false, false, nil); er != nil {
		return er
	}
	if er := b.channel.QueueBind(c.C.Queue, c.C.RoutingKey, c.C.Topic, false, nil); er != nil {
		return er
	}
	b.registered = append(b.registered, c.C)
	return nil
}

func (b *Broker) alreadyRegistered(c QueueConfig) bool {
	found := false
	for _, d := range b.registered {
		if c.Topic == d.Topic && c.Queue == d.Queue && c.RoutingKey == d.RoutingKey {
			found = true
			break
		}
	}

	return found
}

func (b *Broker) Publish(queue string, v interface{}) error {
	for _, r := range b.registry {
		if r.C.Queue == queue && b.alreadyRegistered(r.C) {
			body, e := json.Marshal(v)
			if e != nil {
				return e
			}
			return b.channel.PublishWithContext(context.Background(), r.C.Topic, r.C.Queue, false, false, amqp.Publishing{ContentType: "application/json", Body: body})
		}
	}
	return fmt.Errorf("no handler registered for queue %s", queue)
}

func (b *Broker) consume(r QueueRegistry) {
	var err error
	var messages <-chan amqp.Delivery
	if messages, err = b.channel.Consume(r.C.Queue, "", false, false, false, false, nil); err != nil {
		log.Printf("Error on consuming messge: %s", err.Error())
	}
	go func() {
		for msg := range messages {
			if err := r.H(msg.Body); err != nil {
				_ = msg.Nack(false, false)
				continue
			}
			_ = msg.Ack(false)
		}
	}()
	log.Println("[Consumer] Listening:", r.C.Queue)
}

func (b *Broker) Start() error {
	if err := b.setup(); err != nil {
		return err
	}
	for _, registry := range b.registry {
		if b.alreadyRegistered(registry.C) {
			go b.consume(registry)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for _, registry := range b.registry {
					if !b.alreadyRegistered(registry.C) {
						err := b.register(registry)

						if err != nil {
							log.Printf("Error on registering registry: %s", err.Error())
						} else {
							b.registered = append(b.registered, registry.C)
						}
					}
				}
				if len(b.registered) == len(b.registry) {
					cancel()
				}
			}

		}
	}()

	return nil
}

func NewBroker() *Broker {
	return &Broker{}
}
