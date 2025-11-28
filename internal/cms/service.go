package cms

import (
	"fmt"

	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"gorm.io/gorm"
)

type CmsService interface {
	GetTranslations(language string) []Translation
	GetTranslation(code string, language string) (Translation, error)
}

type cmsService struct {
	cmsRepository commonRepository.CrudRepository[CmsModel]
}

func (s *cmsService) GetTranslations(language string) []Translation {
	c := commonRepository.SearchCriteria{}
	if language != "" {
		c.WhereClause = commonRepository.WhereClause{
			Query:  "language = ?",
			Params: []interface{}{language},
		}
	}
	cms := s.cmsRepository.FindAll(c)
	translations := make([]Translation, len(cms))
	for i, cm := range cms {
		translations[i] = Translation{
			Code:     cm.Code,
			Value:    cm.Value,
			Language: cm.Language,
		}
	}

	return translations
}

func (s *cmsService) GetTranslation(code string, language string) (Translation, error) {
	cms := s.cmsRepository.FindAll(commonRepository.SearchCriteria{WhereClause: commonRepository.WhereClause{
		Query:  "code = ? AND language = ?",
		Params: []interface{}{code, language},
	}})

	if len(cms) == 0 {
		return Translation{}, fmt.Errorf("cms not found")
	}
	cm := cms[0]
	return Translation{
		Code:     cm.Code,
		Value:    cm.Value,
		Language: cm.Language,
	}, nil
}

func NewCmsService(DB *gorm.DB) CmsService {
	cmsRepository := commonRepository.NewCrudRepository[CmsModel](DB)
	return &cmsService{cmsRepository}
}
