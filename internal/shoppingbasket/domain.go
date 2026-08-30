package shoppingbasket

import (
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
)

type ShoppingBasket struct {
	ID            Types.ID
	Items         []ShoppingBasketItem
	ProblemDetail WebHandler.ProblemDetail
}

type ShoppingBasketItem struct {
	Product  Product.Product
	Quantity int
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
