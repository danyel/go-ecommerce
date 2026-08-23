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

//goland:noinspection GoNameStartsWithPackageName
var ShoppingBasketUpdated = Broker.QueueConfig{
	Topic:      ExchangeShoppingBasket,
	Queue:      "shopping-basket.shopping-basket-updated",
	RoutingKey: "shopping-basket-updated",
}

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketEventHandler interface {
	handleShoppingBasketUpdated(body []byte) error
	handleShoppingBasketCreated(body []byte) error
}

type shoppingBasketEvent struct {
	s ShoppingBasketService
}

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketCreatedEvent struct {
	ID Types.ID `json:"id"`
}

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketUpdatedEvent struct {
	ID        Types.ID `json:"id"`
	ProductID Types.ID
	Quantity  int
}

func (shoppingBasketEvent *shoppingBasketEvent) handleShoppingBasketUpdated(body []byte) error {
	var shoppingBasketUpdatedEvent ShoppingBasketUpdatedEvent
	if err := JSON.Unmarshal(body, &shoppingBasketUpdatedEvent); err != nil {
		Log.Printf("Error unmarshalling ShoppingBasketUpdated event: %v\n", err)
		return err
	}
	Log.Printf("Event [ShoppingBasketUpdated] Received: %v", shoppingBasketUpdatedEvent)
	return nil
}

func (shoppingBasketEvent *shoppingBasketEvent) handleShoppingBasketCreated(body []byte) error {
	var shoppingBasketCreatedEvent ShoppingBasketCreatedEvent
	if err := JSON.Unmarshal(body, &shoppingBasketCreatedEvent); err != nil {
		return err
	}
	Log.Printf("Event [ShoppingBasketCreated] Received: %v", shoppingBasketCreatedEvent)
	return nil
}

func RegisterShoppingBasketEvents(s ShoppingBasketService, b *Broker.MessageBroker) {
	h := &shoppingBasketEvent{s}
	b.RegisterConsumer(ShoppingBasketCreated, h.handleShoppingBasketCreated)
	b.RegisterConsumer(ShoppingBasketUpdated, h.handleShoppingBasketUpdated)
}
