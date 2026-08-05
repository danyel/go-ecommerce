package reservation

import "github.com/google/uuid"

type ReservationAggregate struct {
	ProductId        uuid.UUID
	ShoppingBasketId uuid.UUID
	Quantity         int16
}
