package shopping_basket

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
	HandleShoppingBasketUpdated(body []byte) error
	HandleShoppingBasketCreated(body []byte) error
}

type shoppingBasketEvent struct {
	s ShoppingBasketService
}

type ShoppingBasketCreatedEvent struct {
	Id Types.Id `json:"id"`
}

type ShoppingBasketUpdatedEvent struct {
	Id        Types.Id `json:"id"`
	ProductId Types.Id
	Quantity  int
}

func (s *shoppingBasketEvent) HandleShoppingBasketUpdated(body []byte) error {
	var event ShoppingBasketUpdatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		Log.Printf("Error unmarshalling ShoppingBasketUpdated event: %v\n", err)
		return err
	}
	Log.Printf("Event [ShoppingBasketUpdated] Received: %v", event)
	return nil
}

func (s *shoppingBasketEvent) HandleShoppingBasketCreated(body []byte) error {
	var event ShoppingBasketCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Printf("Event [ShoppingBasketCreated] Received: %v", event)
	return nil
}

func NewShoppingBasketEventHandler(s ShoppingBasketService) ShoppingBasketEventHandler {
	h := &shoppingBasketEvent{s}
	return h
}
