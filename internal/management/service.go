package management

import (
	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"gorm.io/gorm"
)

type CmsService interface {
	CreateTranslation(createCms CreateCms) (cmsId CmsId, err error)
}

type cmsService struct {
	cmsRepository commonRepository.CrudRepository[cms.CmsModel]
}

func (s *cmsService) CreateTranslation(createCms CreateCms) (cmsId CmsId, err error) {
	cmsModel := &cms.CmsModel{
		Value:    createCms.Value,
		Language: createCms.Language,
	}
	if err = s.cmsRepository.Create(cmsModel); err != nil {
		return CmsId{}, err
	}
	return CmsId{cmsModel.ID}, nil
}

func NewCmsService(DB *gorm.DB) CmsService {
	return &cmsService{commonRepository.NewCrudRepository[cms.CmsModel](DB)}
}
