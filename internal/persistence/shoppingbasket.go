package persistence

import (
	Time "time"

	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

type ShoppingBasketModel struct {
	ID        Uuid.UUID                 `gorm:"type:uuid;primaryKey"`
	UserID    *Uuid.UUID                `gorm:"type:uuid;index"`
	Status    string                    `gorm:"type:text;not null;default:OPEN"`
	Items     []ShoppingBasketItemModel `gorm:"foreignKey:ShoppingBasketID"`
	CreatedAt Time.Time
	UpdatedAt Time.Time
}

type ShoppingBasketItemModel struct {
	ID               Uuid.UUID `gorm:"type:uuid;primaryKey"`
	ShoppingBasketID Uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID        Uuid.UUID `gorm:"type:uuid;not null;index"`
	Price            float64   `gorm:"column:unit_price;type:numeric(12,2)"`
	Quantity         int
	MaxQuantity      int
}

func (c *ShoppingBasketModel) TableName() string {
	return "ecommerce.shopping_baskets"
}

func (c *ShoppingBasketModel) BeforeCreate(_ *Database.DB) (err error) {
	if c.ID == Uuid.Nil {
		c.ID = Uuid.New()
	}
	return
}

func (c *ShoppingBasketItemModel) TableName() string {
	return "ecommerce.shopping_basket_items"
}

func (c *ShoppingBasketItemModel) BeforeCreate(_ *Database.DB) (err error) {
	if c.ID == Uuid.Nil {
		c.ID = Uuid.New()
	}
	return
}
