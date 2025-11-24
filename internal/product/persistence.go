package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Brand       string
	Name        string
	Description string
	Code        string
	Price       uint32
	Category    uuid.UUID
	ImageUrl    string
	Stock       uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *ProductModel) TableName() string {
	return "ecommerce.products"
}

func (p *ProductModel) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return
}
