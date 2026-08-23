package product

import (
	Time "time"

	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductModel struct {
	ID          Uuid.UUID `gorm:"type:uuid;primaryKey"`
	Brand       string
	Name        string
	Description string
	Code        string
	Price       float64 `gorm:"type:numeric(10,2)"`
	CategoryID  Uuid.UUID
	ImageURL    string
	Stock       int
	CreatedAt   Time.Time
	UpdatedAt   Time.Time
}

func (productModel *ProductModel) TableName() string {
	return "products"
}

func (productModel *ProductModel) BeforeCreate(_ *Database.DB) (err error) {
	if productModel.ID == Uuid.Nil {
		productModel.ID = Uuid.New()
	}

	return
}
