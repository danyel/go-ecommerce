package shopping_basket

import (
	"time"

	"github.com/danyel/ecommerce/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShoppingBasketModel struct {
	ID        uuid.UUID               `gorm:"type:uuid;primaryKey"`
	Items     []*product.ProductModel `gorm:"many2many:shopping_basket_items;joinForeignKey:ID;joinForeignKey:ShoppingBasketID;joinReferences:ID;joinReferences:ProductID"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
