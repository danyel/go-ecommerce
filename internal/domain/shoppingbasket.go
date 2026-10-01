package domain

import (
	Types "github.com/danyel/ecommerce/internal/types"
	Uuid "github.com/google/uuid"
)

type ShoppingBasket struct {
	ID            Types.ID
	Items         []ShoppingBasketItem
	ProblemDetail ProblemDetail
}

type ShoppingBasketItem struct {
	Product  Product
	Quantity int
}

type ShoppingBasketItemUpdate struct {
	ProductID Uuid.UUID
	Quantity  int
}

// StockAdjustment returns the change to available stock when a basket line
// changes from currentQuantity to requestedQuantity.
func StockAdjustment(currentQuantity, requestedQuantity int) int {
	return currentQuantity - requestedQuantity
}

func EmptyShoppingBasket() ShoppingBasket {
	return ShoppingBasket{}
}

func (shoppingBasket *ShoppingBasket) TotalPrice() Types.Price {
	var totalPrice = Types.Price{
		Tax:       Types.Float64(0),
		Inclusive: Types.Float64(0),
		Exclusive: Types.Float64(0),
		Currency:  "EUR",
	}
	for _, item := range shoppingBasket.Items {
		calculatePrice := item.CalculatePrice(item.Quantity)
		totalPrice.Tax += calculatePrice.Tax
		totalPrice.Inclusive += calculatePrice.Inclusive
		totalPrice.Exclusive += calculatePrice.Exclusive
	}

	return totalPrice
}

func (shoppingBasketItem *ShoppingBasketItem) CalculatePrice(quantity int) Types.Price {
	return Types.Price{
		Currency:  shoppingBasketItem.Product.Price.Currency,
		Inclusive: shoppingBasketItem.Product.Price.Inclusive * Types.Float64(quantity),
		Tax:       shoppingBasketItem.Product.Price.Tax * Types.Float64(quantity),
		Exclusive: shoppingBasketItem.Product.Price.Exclusive * Types.Float64(quantity),
	}
}
