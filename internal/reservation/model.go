package reservation

import (
	Types "github.com/danyel/ecommerce/internal/common/types"
)

type CreateReservation struct {
	ShoppingBasketID Types.ID `json:"shopping_basket_id"`
	ProductID        Types.ID `json:"product_id"`
	Quantity         int      `json:"quantity"`
}

type Reservation struct {
	ShoppingBasketID Types.ID `json:"shopping_basket_id"`
	ProductID        Types.ID `json:"product_id"`
	Quantity         int      `json:"quantity"`
}
