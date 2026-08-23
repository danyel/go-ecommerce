package management

import (
	CMS "github.com/danyel/ecommerce/internal/cms"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementService interface {
	CreateTranslation(createCms CreateCms) (ID CmsID, err error)
}

type managementService struct {
	cmsRepository Repository.CrudRepository[CMS.CmsModel]
}

func (managementService *managementService) CreateTranslation(createCms CreateCms) (ID CmsID, err error) {
	cmsModel := &CMS.CmsModel{
		Code:     createCms.Code,
		Value:    createCms.Value,
		Language: createCms.Language,
	}
	if err = managementService.cmsRepository.Create(cmsModel); err != nil {
		return CmsID{}, err
	}
	return CmsID{cmsModel.ID}, nil
}

func NewService(cmsRepository Repository.CrudRepository[CMS.CmsModel]) ManagementService {
	return &managementService{cmsRepository: cmsRepository}
}
