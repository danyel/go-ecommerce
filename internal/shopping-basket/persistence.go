package shopping_basket

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShoppingBasketModel struct {
	ID        uuid.UUID                 `gorm:"type:uuid;primaryKey"`
	Items     []ShoppingBasketItemModel `gorm:"foreignKey:ShoppingBasketID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ShoppingBasketItemModel struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	ShoppingBasketID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductId        uuid.UUID `gorm:"type:uuid;not null;index"`
	Price            int
	Quantity         int
}

func (c *ShoppingBasketModel) TableName() string {
	return "shopping_basket"
}

func (c *ShoppingBasketModel) BeforeCreate(_ *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

func (c *ShoppingBasketItemModel) TableName() string {
	return "shopping_basket_items"
}

func (c *ShoppingBasketItemModel) BeforeCreate(_ *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
