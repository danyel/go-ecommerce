package shopping_basket

import (
	Types "github.com/danyel/ecommerce/internal/common/types"
	Uuid "github.com/google/uuid"
)

type UpdateShoppingBasketItem struct {
	ProductId Types.Id `json:"product_id"`
	Quantity  int      `json:"quantity"`
}

type ShoppingBasketId struct {
	ID Types.Id `json:"id"`
}

func NewShoppingBasketId(id Uuid.UUID) Types.Id {
	return Types.Id{ID: id}
}

func EmptyShoppingBasket() ShoppingBasket {
	return ShoppingBasket{}
}

type ShoppingBasketItem struct {
	ID         Types.Id    `json:"id"`
	Name       string      `json:"name"`
	BasePrice  Types.Price `json:"base_price"`
	TotalPrice Types.Price `json:"total_price"`
	ProductId  Types.Id    `json:"product_id"`
	ImageUrl   string      `json:"image_url"`
	Quantity   int         `json:"quantity"`
}

type ShoppingBasket struct {
	ID         Types.Id             `json:"id"`
	Items      []ShoppingBasketItem `json:"items"`
	TotalPrice Types.Price          `json:"total_price"`
}
