package cms

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsModel struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;"`
	Code      string    `gorm:"type:text;"`
	Value     string    `gorm:"type:text;"`
	Language  string    `gorm:"type:varchar(5);"`
	CreatedAt time.Time `gorm:"type:timestamp;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (p *CmsModel) TableName() string {
	return "cms"
}

func (p *CmsModel) BeforeCreate(_ *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return
}
