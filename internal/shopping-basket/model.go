package shopping_basket

import (
	productmanagement "github.com/danyel/ecommerce/internal/product-management"
	"github.com/google/uuid"
)

type AddItem struct {
	ProductId uuid.UUID `json:"product_id"`
}

type ShoppingId struct {
	Id uuid.UUID `json:"id"`
}

type ShoppingBasket struct {
	Id                  uuid.UUID                   `json:"id"`
	Items               []productmanagement.Product `json:"items"`
	TotalPriceInclusive float32                     `json:"total_price_inclusive"`
	Tax                 float32                     `json:"tax"`
	TotalPriceExclusive float32                     `json:"total_price_exclusive"`
}
