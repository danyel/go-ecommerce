package category

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryModel struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Name      string           `gorm:"type:text;not null"`
	Children  []*CategoryModel `gorm:"many2many:category_children;joinForeignKey:ParentID;joinReferences:ChildID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *CategoryModel) TableName() string {
	return "ecommerce.categories"
}

func (c *CategoryModel) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
