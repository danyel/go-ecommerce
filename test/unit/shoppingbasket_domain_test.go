package unit

import (
	UnitTesting "testing"

	Domain "github.com/danyel/ecommerce/internal/domain"
	Types "github.com/danyel/ecommerce/internal/types"
	Require "github.com/stretchr/testify/require"
)

func TestShoppingBasketStockAdjustment(t *UnitTesting.T) {
	Require.Equal(t, 0, Domain.StockAdjustment(0, 0))
	Require.Equal(t, -3, Domain.StockAdjustment(0, 3))
	Require.Equal(t, 2, Domain.StockAdjustment(3, 1))
	Require.Equal(t, -4, Domain.StockAdjustment(1, 5))
}

func TestShoppingBasketTotalPrice(t *UnitTesting.T) {
	basket := Domain.ShoppingBasket{Items: []Domain.ShoppingBasketItem{
		{
			Product:  Domain.Product{Price: Types.NewPrice(10, "EUR")},
			Quantity: 2,
		},
		{
			Product:  Domain.Product{Price: Types.NewPrice(5, "EUR")},
			Quantity: 1,
		},
	}}

	Require.Equal(t, Types.Float64(25), basket.TotalPrice().Inclusive)
	Require.Equal(t, "EUR", basket.TotalPrice().Currency)
}
