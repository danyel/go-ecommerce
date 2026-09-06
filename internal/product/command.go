package product

import (
	Uuid "github.com/google/uuid"
)

type UpdateProductStockCommand struct {
	ShoppingBasketId Uuid.UUID `json:"shopping_basket_id"`
	ProductID        Uuid.UUID `json:"product_id"`
	Quantity         int       `json:"quantity"` // can be positive or negative
}
