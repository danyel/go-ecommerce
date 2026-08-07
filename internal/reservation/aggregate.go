package reservation

import Uuid "github.com/google/uuid"

//goland:noinspection GoNameStartsWithPackageName
type ReservationAggregate struct {
	ProductId        Uuid.UUID
	ShoppingBasketId Uuid.UUID
	Quantity         int16
}
