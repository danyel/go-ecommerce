package reservation

import (
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationModel struct {
	ShoppingBasketId Uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductId        Uuid.UUID `gorm:"type:uuid;not null"`
	Quantity         int
}

func (c *ReservationModel) TableName() string {
	return "reservations"
}
