package service

import (
	Fmt "fmt"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Model "github.com/danyel/ecommerce/internal/model"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
)

//goland:noinspection GoNameStartsWithPackageName
type ICmsService interface {
	GetTranslations(language string) []Model.Translation
	GetTranslation(code string, language string) (Model.Translation, error)
}

type cmsService struct {
	cmsRepository Persistence.CrudRepository[Persistence.CmsModel]
}

func (cmsService *cmsService) GetTranslations(language string) []Model.Translation {
	searchCriteria := Persistence.SearchCriteria{}
	if language != "" {
		searchCriteria.WhereClause = Persistence.WhereClause{
			Query:  "language = ?",
			Params: []any{language},
		}
	}
	translation := cmsService.cmsRepository.FindAll(searchCriteria)
	translations := make([]Model.Translation, len(translation))
	for index, cms := range translation {
		translations[index] = Model.Translation{
			Code:     cms.Code,
			Value:    cms.Value,
			Language: cms.Language,
		}
	}
	Logger.Log.Debug("GetTranslations: %v", translations)
	return translations
}

func (cmsService *cmsService) GetTranslation(code string, language string) (Model.Translation, error) {
	translations := cmsService.cmsRepository.FindAll(Persistence.SearchCriteria{WhereClause: Persistence.WhereClause{
		Query:  "code = ? AND language = ?",
		Params: []any{code, language},
	}})

	if len(translations) == 0 {
		return Model.Translation{}, Fmt.Errorf("cms not found")
	}
	cms := translations[0]
	return Model.Translation{
		Code:     cms.Code,
		Value:    cms.Value,
		Language: cms.Language,
	}, nil
}

func CmsService(cmsRepository Persistence.CrudRepository[Persistence.CmsModel]) ICmsService {
	return &cmsService{cmsRepository}
}
