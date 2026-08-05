package shopping_basket

import (
	"encoding/json"

	"github.com/danyel/ecommerce/cmd/broker"
	"github.com/danyel/ecommerce/internal/common/types"
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
	Id types.Id `json:"id"`
}

type ShoppingBasketUpdatedEvent struct {
	Id        types.Id `json:"id"`
	ProductId types.Id
	Quantity  int
}

func HandleShoppingBasketUpdated(body []byte) error {
	var event ShoppingBasketUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	return nil
}

func HandleShoppingBasketCreated(body []byte) error {
	var event ShoppingBasketCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	return nil
}
