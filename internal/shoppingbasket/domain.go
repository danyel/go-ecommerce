package shoppingbasket

import Types "github.com/danyel/ecommerce/internal/common/types"

type ShoppingBasketDTO struct {
	ID    Types.ID
	Items []ShoppingBasketItemDTO
}

type ShoppingBasketItemDTO struct {
	ProductID Types.ID
	Quantity  int
}
