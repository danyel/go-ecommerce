package reservation

import (
	Types "github.com/danyel/ecommerce/internal/common/types"
)

type CreateReservation struct {
	ShoppingBasketId Types.Id `json:"shopping_basket_id"`
	ProductId        Types.Id `json:"product_id"`
	Quantity         int      `json:"quantity"`
}

type Reservation struct {
	ShoppingBasketId Types.Id `json:"shopping_basket_id"`
	ProductId        Types.Id `json:"product_id"`
	Quantity         int      `json:"quantity"`
}
