package shopping_basket

import (
	"github.com/google/uuid"
)

type AddItem struct {
	ProductId uuid.UUID `json:"product_id"`
}

type ShoppingId struct {
	ID uuid.UUID `json:"id"`
}

type ShoppingBasketItem struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Price    int       `json:"price"`
	ImageUrl string    `json:"image_url"`
	Amount   int       `json:"amount"`
}

type ShoppingBasket struct {
	ID                  uuid.UUID            `json:"id"`
	Items               []ShoppingBasketItem `json:"items"`
	TotalPriceInclusive float32              `json:"total_price_inclusive"`
	Tax                 float32              `json:"tax"`
	TotalPriceExclusive float32              `json:"total_price_exclusive"`
}
