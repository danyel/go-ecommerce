package reservation

import (
	"github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationModel struct {
	ShoppingBasketId uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductId        uuid.UUID `gorm:"type:uuid;not null"`
	Quantity         int
}

func (c *ReservationModel) TableName() string {
	return "reservations"
}
