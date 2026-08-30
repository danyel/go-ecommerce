package shoppingbasket

import (
	Types "github.com/danyel/ecommerce/internal/common/types"
)

type UpdateShoppingBasketItemDTO struct {
	ProductID Types.ID `json:"product_id"`
	Quantity  int      `json:"quantity"`
}

type ShoppingBasketItemDTO struct {
	Name       string         `json:"name"`
	BasePrice  Types.PriceDTO `json:"base_price"`
	TotalPrice Types.PriceDTO `json:"total_price"`
	ProductID  Types.ID       `json:"product_id"`
	ImageURL   string         `json:"image_url"`
	Quantity   int            `json:"quantity"`
}

type PromoDTO struct {
	Code       string        `json:"code"`
	Percentage Types.Float64 `json:"percentage"`
}

type DiscountDTO struct{}

type ShoppingBasketDTO struct {
	ID         Types.ID                `json:"id"`
	Items      []ShoppingBasketItemDTO `json:"items"`
	TotalPrice Types.PriceDTO          `json:"total_price"`
	Discounts  []DiscountDTO           `json:"discounts"`
	Promos     []PromoDTO              `json:"promos"`
}
