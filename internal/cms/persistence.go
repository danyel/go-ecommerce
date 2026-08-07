package cms

import (
	Time "time"

	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsModel struct {
	ID        Uuid.UUID `gorm:"primary_key;type:uuid;"`
	Code      string    `gorm:"type:text;"`
	Value     string    `gorm:"type:text;"`
	Language  string    `gorm:"type:varchar(5);"`
	CreatedAt Time.Time `gorm:"type:timestamp;"`
	UpdatedAt Time.Time `gorm:"type:timestamp;"`
}

func (p *CmsModel) TableName() string {
	return "cms"
}

func (p *CmsModel) BeforeCreate(_ *Database.DB) (err error) {
	if p.ID == Uuid.Nil {
		p.ID = Uuid.New()
	}

	return
}
