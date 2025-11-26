package management

import (
	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"gorm.io/gorm"
)

type ManagementService interface {
	CreateTranslation(createCms CreateCms) (cmsId CmsId, err error)
}

type managementService struct {
	cmsRepository commonRepository.CrudRepository[cms.CmsModel]
}

func (s *managementService) CreateTranslation(createCms CreateCms) (cmsId CmsId, err error) {
	cmsModel := &cms.CmsModel{
		Value:    createCms.Value,
		Language: createCms.Language,
	}
	if err = s.cmsRepository.Create(cmsModel); err != nil {
		return CmsId{}, err
	}
	return CmsId{cmsModel.ID}, nil
}

func NewManagementService(DB *gorm.DB) ManagementService {
	return &managementService{commonRepository.NewCrudRepository[cms.CmsModel](DB)}
}
