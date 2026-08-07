package shopping_basket

import (
	Types "github.com/danyel/ecommerce/internal/common/types"
)

type UpdateShoppingBasketItem struct {
	ProductId Types.Id `json:"product_id"`
	Quantity  int      `json:"quantity"`
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

type Promo struct {
	Code       string        `json:"code"`
	Percentage Types.Float64 `json:"percentage"`
}

type Discount struct {
}

type ShoppingBasket struct {
	ID         Types.Id             `json:"id"`
	Items      []ShoppingBasketItem `json:"items"`
	TotalPrice Types.Price          `json:"total_price"`
	Discounts  []Discount           `json:"discounts"`
	Promos     []Promo              `json:"promos"`
}
