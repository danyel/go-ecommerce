package reservation

import Uuid "github.com/google/uuid"

type ReservationAggregate struct {
	ProductId        Uuid.UUID
	ShoppingBasketId Uuid.UUID
	Quantity         int16
}
