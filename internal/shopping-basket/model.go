package shopping_basket

import (
	"github.com/google/uuid"
)

type UpdateShoppingBasketItem struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type ShoppingId struct {
	ID uuid.UUID `json:"id"`
}

type ShoppingBasketItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	ProductId uuid.UUID `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	Quantity  int       `json:"quantity"`
}

type ShoppingBasket struct {
	ID                  uuid.UUID            `json:"id"`
	Items               []ShoppingBasketItem `json:"items"`
	TotalPriceInclusive float32              `json:"total_price_inclusive"`
	Tax                 float32              `json:"tax"`
	TotalPriceExclusive float32              `json:"total_price_exclusive"`
}
