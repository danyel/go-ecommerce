package category

import (
	"encoding/json"
	"log"

	"github.com/danyel/ecommerce/cmd/broker"
)

const (
	ExchangeCategory     = "category.topic"
	QueueCategoryCreated = "categories.category_created"
)

var CategoryCreated = broker.QueueConfig{
	Topic: ExchangeCategory,
	Queue: QueueCategoryCreated,
}

var CategoryCreated2 = broker.QueueConfig{
	Topic: ExchangeCategory,
	Queue: "categories.category_created6",
}

type CategoryCreatedEvent struct {
	Id string `json:"id"`
}

func HandleCategoryCreated2(body []byte) error {
	var event CategoryCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	log.Println(event.Id)
	return nil
}
func HandleCategoryCreated(body []byte) error {
	var event CategoryCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	log.Println(event.Id)
	return nil
}
