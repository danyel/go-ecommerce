package shoppingbasket

import (
	Types "github.com/danyel/ecommerce/internal/common/types"
)

type UpdateShoppingBasketItem struct {
	ProductID Types.ID `json:"product_id"`
	Quantity  int      `json:"quantity"`
}

func EmptyShoppingBasket() ShoppingBasket {
	return ShoppingBasket{}
}

type ShoppingBasketItem struct {
	ID         Types.ID    `json:"id"`
	Name       string      `json:"name"`
	BasePrice  Types.Price `json:"base_price"`
	TotalPrice Types.Price `json:"total_price"`
	ProductID  Types.ID    `json:"product_id"`
	ImageURL   string      `json:"image_url"`
	Quantity   int         `json:"quantity"`
}

type Promo struct {
	Code       string        `json:"code"`
	Percentage Types.Float64 `json:"percentage"`
}

type Discount struct{}

type ShoppingBasket struct {
	ID         Types.ID             `json:"id"`
	Items      []ShoppingBasketItem `json:"items"`
	TotalPrice Types.Price          `json:"total_price"`
	Discounts  []Discount           `json:"discounts"`
	Promos     []Promo              `json:"promos"`
}
