package category

import (
	Time "time"

	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryModel struct {
	ID        Uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Name      string           `gorm:"type:text;not null"`
	Children  []*CategoryModel `gorm:"many2many:category_children;joinForeignKey:ParentID;joinReferences:ChildID"`
	CreatedAt Time.Time
	UpdatedAt Time.Time
}

func (c *CategoryModel) TableName() string {
	return "categories"
}

func (c *CategoryModel) BeforeCreate(_ *Database.DB) (err error) {
	if c.ID == Uuid.Nil {
		c.ID = Uuid.New()
	}
	return
}
