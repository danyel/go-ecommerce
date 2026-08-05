package reservation

import "github.com/google/uuid"

type CreateReservation struct {
	ShoppingBasketId uuid.UUID `json:"shopping_basket_id"`
	ProductId        uuid.UUID `json:"product_id"`
	Quantity         int       `json:"quantity"`
}

type Reservation struct {
	ShoppingBasketId uuid.UUID `json:"shopping_basket_id"`
	ProductId        uuid.UUID `json:"product_id"`
	Quantity         int       `json:"quantity"`
}

//goland:noinspection GoNameStartsWithPackageName
type ReservationId struct {
	ID uuid.UUID `json:"id"`
}
