package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/danyel/ecommerce/cmd/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueConfig struct {
	Topic      string
	Queue      string
	RoutingKey string
}

type Broker struct {
	registry []QueueRegistry
	channel  *amqp.Channel
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
		if er := b.channel.ExchangeDeclare(c.C.Topic, "topic", true, false, false, false, nil); er != nil {
			return er
		}
		if _, er := b.channel.QueueDeclare(c.C.Queue, true, false, false, false, nil); er != nil {
			return er
		}
		if er := b.channel.QueueBind(c.C.Queue, c.C.RoutingKey, c.C.Topic, false, nil); er != nil {
			return er
		}
	}

	return nil
}

func (b *Broker) Publish(queue string, v interface{}) error {
	for _, r := range b.registry {
		if r.C.Queue == queue {
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
		log.Printf(err.Error())
	}

	for {
		if messages, err = b.channel.Consume(r.C.Queue, "", false, false, false, false, nil); err != nil {
			log.Printf(err.Error())
			time.Sleep(30 * time.Second)
		} else {
			break
		}
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
		go b.consume(registry)
	}

	return nil
}

func NewBroker() *Broker {
	return &Broker{}
}
