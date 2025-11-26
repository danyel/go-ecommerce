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
	CategoryId  uuid.UUID `gorm:"type:uuid"`
	ImageUrl    string
	Stock       uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *ProductModel) TableName() string {
	return "products"
}

func (p *ProductModel) BeforeCreate(_ *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return
}
