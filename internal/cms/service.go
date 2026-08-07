package cms

import (
	Fmt "fmt"

	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsService interface {
	GetTranslations(language string) []Translation
	GetTranslation(code string, language string) (Translation, error)
}

type cmsService struct {
	cmsRepository CommonRepository.CrudRepository[CmsModel]
}

func (s *cmsService) GetTranslations(language string) []Translation {
	c := CommonRepository.SearchCriteria{}
	if language != "" {
		c.WhereClause = CommonRepository.WhereClause{
			Query:  "language = ?",
			Params: []any{language},
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
	cms := s.cmsRepository.FindAll(CommonRepository.SearchCriteria{WhereClause: CommonRepository.WhereClause{
		Query:  "code = ? AND language = ?",
		Params: []interface{}{code, language},
	}})

	if len(cms) == 0 {
		return Translation{}, Fmt.Errorf("cms not found")
	}
	cm := cms[0]
	return Translation{
		Code:     cm.Code,
		Value:    cm.Value,
		Language: cm.Language,
	}, nil
}

func NewCmsService(DB *Database.DB) CmsService {
	cmsRepository := CommonRepository.NewCrudRepository[CmsModel](DB)
	return &cmsService{cmsRepository}
}
