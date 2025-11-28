package shopping_basket

import (
	"github.com/danyel/ecommerce/internal/product"
	"github.com/google/uuid"
)

type AddItem struct {
	ProductId uuid.UUID `json:"product_id"`
}

type ShoppingId struct {
	Id uuid.UUID `json:"id"`
}

type ShoppingBasket struct {
	Id                  uuid.UUID         `json:"id"`
	Items               []product.Product `json:"items"`
	TotalPriceInclusive float32           `json:"total_price_inclusive"`
	Tax                 float32           `json:"tax"`
	TotalPriceExclusive float32           `json:"total_price_exclusive"`
}
