package shopping_basket

import (
	"encoding/json"
	"log"

	"github.com/danyel/ecommerce/cmd/broker"
	"github.com/google/uuid"
)

const (
	ExchangeShoppingBasket          = "shopping-basket.topic"
	ShoppingBasketCreatedRoutingKey = "shopping-basket-created"
	ShoppingBasketUpdatedRoutingKey = "shopping-basket-updated"
	ShoppingBasketCreatedQueue      = "shopping-basket.shopping-basket-created"
	ShoppingBasketUpdatedQueue      = "shopping-basket.shopping-basket-updated"
)

var ShoppingBasketCreated = broker.QueueConfig{
	Topic:      ExchangeShoppingBasket,
	Queue:      ShoppingBasketCreatedQueue,
	RoutingKey: ShoppingBasketCreatedRoutingKey,
}

var ShoppingBasketUpdated = broker.QueueConfig{
	Topic:      ExchangeShoppingBasket,
	Queue:      ShoppingBasketUpdatedQueue,
	RoutingKey: ShoppingBasketUpdatedRoutingKey,
}

type ShoppingBasketCreatedEvent struct {
	Id ShoppingBasketId `json:"id"`
}

type ShoppingBasketUpdatedEvent struct {
	Id        ShoppingBasketId `json:"id"`
	ProductId uuid.UUID
	Quantity  int
}

func HandleShoppingBasketUpdated(body []byte) error {
	var event ShoppingBasketUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	log.Println(event)
	return nil
}

func HandleShoppingBasketCreated(body []byte) error {
	var event ShoppingBasketCreatedEvent
	log.Printf("Entering Shopping Basket Created")
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	log.Println(event.Id)
	return nil
}
