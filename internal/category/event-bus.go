// Package category: events
// author: Daniel Noulet
package category

import (
	JSON "encoding/json"
	Log "log"

	Broker "github.com/danyel/ecommerce/cmd/broker"
)

const (
	ExchangeCategory     = "category.topic"
	QueueCategoryCreated = "categories.category_created"
)

//goland:noinspection GoUnusedGlobalVariable,GoNameStartsWithPackageName
var CategoryCreated = Broker.QueueConfig{
	Topic: ExchangeCategory,
	Queue: QueueCategoryCreated,
}

//goland:noinspection GoNameStartsWithPackageName
type CategoryCreatedEvent struct {
	ID string `json:"id"`
}

//goland:noinspection GoUnusedExportedFunction
func HandleCategoryCreated2(body []byte) error {
	var event CategoryCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Println(event.ID)
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func HandleCategoryCreated(body []byte) error {
	var event CategoryCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Println(event.ID)
	return nil
}
