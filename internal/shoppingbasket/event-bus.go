package shoppingbasket

import (
	JSON "encoding/json"
	Log "log"

	Broker "github.com/danyel/ecommerce/cmd/broker"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

const (
	ExchangeShoppingBasket = "shopping-basket.topic"
)

var ShoppingBasketCreated = Broker.QueueConfig{
	Topic:      ExchangeShoppingBasket,
	Queue:      "shopping-basket.shopping-basket-created",
	RoutingKey: "shopping-basket-created",
}

var ShoppingBasketUpdated = Broker.QueueConfig{
	Topic:      ExchangeShoppingBasket,
	Queue:      "shopping-basket.shopping-basket-updated",
	RoutingKey: "shopping-basket-updated",
}

type ShoppingBasketEventHandler interface {
	handleShoppingBasketUpdated(body []byte) error
	handleShoppingBasketCreated(body []byte) error
}

type shoppingBasketEvent struct {
	s ShoppingBasketService
}

type ShoppingBasketCreatedEvent struct {
	ID Types.ID `json:"id"`
}

type ShoppingBasketUpdatedEvent struct {
	ID        Types.ID `json:"id"`
	ProductID Types.ID
	Quantity  int
}

func (s *shoppingBasketEvent) handleShoppingBasketUpdated(body []byte) error {
	var event ShoppingBasketUpdatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		Log.Printf("Error unmarshalling ShoppingBasketUpdated event: %v\n", err)
		return err
	}
	Log.Printf("Event [ShoppingBasketUpdated] Received: %v", event)
	return nil
}

func (s *shoppingBasketEvent) handleShoppingBasketCreated(body []byte) error {
	var event ShoppingBasketCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Printf("Event [ShoppingBasketCreated] Received: %v", event)
	return nil
}

func RegisterShoppingBasketEvents(s ShoppingBasketService, b *Broker.Broker) {
	h := &shoppingBasketEvent{s}
	b.RegisterConsumer(ShoppingBasketCreated, h.handleShoppingBasketCreated)
	b.RegisterConsumer(ShoppingBasketUpdated, h.handleShoppingBasketUpdated)
}
