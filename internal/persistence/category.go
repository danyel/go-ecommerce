package persistence

import (
	Time "time"

	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryModel struct {
	ID        Uuid.UUID        `gorm:"type:uuid;primaryKey"`
	ParentID  *Uuid.UUID       `gorm:"type:uuid;index"`
	Name      string           `gorm:"type:text;not null"`
	Slug      string           `gorm:"type:text;not null"`
	Children  []*CategoryModel `gorm:"foreignKey:ParentID"`
	CreatedAt Time.Time
	UpdatedAt Time.Time
}

func (categoryModel *CategoryModel) TableName() string {
	return "ecommerce.categories"
}

func (categoryModel *CategoryModel) BeforeCreate(_ *Database.DB) (err error) {
	if categoryModel.ID == Uuid.Nil {
		categoryModel.ID = Uuid.New()
	}
	if categoryModel.Slug == "" {
		categoryModel.Slug = categoryModel.Name
	}
	return
}
