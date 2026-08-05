package reservation

import (
	"github.com/danyel/ecommerce/internal/common/types"
	"github.com/google/uuid"
)

type CreateReservation struct {
	ShoppingBasketId types.Id `json:"shopping_basket_id"`
	ProductId        types.Id `json:"product_id"`
	Quantity         int      `json:"quantity"`
}

type Reservation struct {
	ShoppingBasketId types.Id `json:"shopping_basket_id"`
	ProductId        types.Id `json:"product_id"`
	Quantity         int      `json:"quantity"`
}

//goland:noinspection GoNameStartsWithPackageName
type ReservationId struct {
	ID uuid.UUID `json:"id"`
}
