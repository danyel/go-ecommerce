// Package category: events
// author: Daniel Noulet
package category

import (
	JSON "encoding/json"

	Broker "github.com/danyel/ecommerce/cmd/broker"
	Logger "github.com/danyel/ecommerce/cmd/logger"
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
	var categoryCreatedEvent CategoryCreatedEvent
	if err := JSON.Unmarshal(body, &categoryCreatedEvent); err != nil {
		return err
	}
	Logger.Log.Debug(categoryCreatedEvent.ID)
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func HandleCategoryCreated(body []byte) error {
	var categoryCreatedEvent CategoryCreatedEvent
	if err := JSON.Unmarshal(body, &categoryCreatedEvent); err != nil {
		return err
	}
	Logger.Log.Debug(categoryCreatedEvent.ID)
	return nil
}
