package reservation

import Uuid "github.com/google/uuid"

//goland:noinspection GoNameStartsWithPackageName
type ReservationAggregate struct {
	ProductID        Uuid.UUID
	ShoppingBasketID Uuid.UUID
	Quantity         int16
}
