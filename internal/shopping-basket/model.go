package shopping_basket

import (
	"github.com/danyel/ecommerce/internal/common/types"
	"github.com/google/uuid"
)

type UpdateShoppingBasketItem struct {
	ProductId types.Id `json:"product_id"`
	Quantity  int      `json:"quantity"`
}

type ShoppingBasketId struct {
	ID types.Id `json:"id"`
}

func NewShoppingBasketId(id uuid.UUID) types.Id {
	return types.Id{ID: id}
}

func EmptyShoppingBasket() ShoppingBasket {
	return ShoppingBasket{}
}

type ShoppingBasketItem struct {
	ID                 types.Id      `json:"id"`
	Name               string        `json:"name"`
	BasePriceInclusive types.Float64 `json:"base_price_inclusive"`
	BasePriceExclusive types.Float64 `json:"base_price_exclusive"`
	BaseTax            types.Float64 `json:"base_tax"`
	PriceInclusive     types.Float64 `json:"price_inclusive"`
	PriceExclusive     types.Float64 `json:"price_exclusive"`
	Tax                types.Float64 `json:"tax"`
	ProductId          types.Id      `json:"product_id"`
	ImageUrl           string        `json:"image_url"`
	Quantity           int           `json:"quantity"`
}

type ShoppingBasket struct {
	ID                  types.Id             `json:"id"`
	Items               []ShoppingBasketItem `json:"items"`
	TotalPriceInclusive types.Float64        `json:"total_price_inclusive"`
	Tax                 types.Float64        `json:"tax"`
	TotalPriceExclusive types.Float64        `json:"total_price_exclusive"`
}
