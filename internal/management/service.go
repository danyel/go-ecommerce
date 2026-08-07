package management

import (
	CMS "github.com/danyel/ecommerce/internal/cms"
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementService interface {
	CreateTranslation(createCms CreateCms) (cmsId CmsId, err error)
}

type managementService struct {
	cmsRepository CommonRepository.CrudRepository[CMS.CmsModel]
}

func (s *managementService) CreateTranslation(createCms CreateCms) (cmsId CmsId, err error) {
	cmsModel := &CMS.CmsModel{
		Code:     createCms.Code,
		Value:    createCms.Value,
		Language: createCms.Language,
	}
	if err = s.cmsRepository.Create(cmsModel); err != nil {
		return CmsId{}, err
	}
	return CmsId{cmsModel.ID}, nil
}

func NewManagementService(DB *Database.DB) ManagementService {
	return &managementService{CommonRepository.NewCrudRepository[CMS.CmsModel](DB)}
}
