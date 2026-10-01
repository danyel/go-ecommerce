package persistence

import (
	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
	Time "time"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationModel struct {
	ID               Uuid.UUID `gorm:"type:uuid;primaryKey"`
	ShoppingBasketID Uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID        Uuid.UUID `gorm:"type:uuid;not null"`
	Quantity         int       `gorm:"not null"`
	CreatedAt        Time.Time
	UpdatedAt        Time.Time
}

func (c *ReservationModel) BeforeCreate(_ *Database.DB) error {
	if c.ID == Uuid.Nil {
		c.ID = Uuid.New()
	}
	return nil
}

func (c *ReservationModel) TableName() string {
	return "ecommerce.reservations"
}
