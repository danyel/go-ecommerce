package cms

import (
	Fmt "fmt"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsService interface {
	GetTranslations(language string) []Translation
	GetTranslation(code string, language string) (Translation, error)
}

type cmsService struct {
	cmsRepository Repository.CrudRepository[CmsModel]
}

func (cmsService *cmsService) GetTranslations(language string) []Translation {
	searchCriteria := Repository.SearchCriteria{}
	if language != "" {
		searchCriteria.WhereClause = Repository.WhereClause{
			Query:  "language = ?",
			Params: []any{language},
		}
	}
	translation := cmsService.cmsRepository.FindAll(searchCriteria)
	translations := make([]Translation, len(translation))
	for index, cms := range translation {
		translations[index] = Translation{
			Code:     cms.Code,
			Value:    cms.Value,
			Language: cms.Language,
		}
	}
	Logger.Log.Debug("GetTranslations: %v", translations)
	return translations
}

func (cmsService *cmsService) GetTranslation(code string, language string) (Translation, error) {
	translations := cmsService.cmsRepository.FindAll(Repository.SearchCriteria{WhereClause: Repository.WhereClause{
		Query:  "code = ? AND language = ?",
		Params: []any{code, language},
	}})

	if len(translations) == 0 {
		return Translation{}, Fmt.Errorf("cms not found")
	}
	cms := translations[0]
	return Translation{
		Code:     cms.Code,
		Value:    cms.Value,
		Language: cms.Language,
	}, nil
}

func NewService(cmsRepository Repository.CrudRepository[CmsModel]) CmsService {
	return &cmsService{cmsRepository}
}
