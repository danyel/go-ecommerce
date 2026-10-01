package service

import (
	Model "github.com/danyel/ecommerce/internal/model"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
)

//goland:noinspection GoNameStartsWithPackageName
type IManagementService interface {
	CreateTranslation(createCms Model.CreateCms) (ID Model.CmsID, err error)
}

type managementService struct {
	cmsRepository Persistence.CrudRepository[Persistence.CmsModel]
}

func (managementService *managementService) CreateTranslation(createCms Model.CreateCms) (ID Model.CmsID, err error) {
	cmsModel := &Persistence.CmsModel{
		Code:     createCms.Code,
		Value:    createCms.Value,
		Language: createCms.Language,
	}
	if err = managementService.cmsRepository.Create(cmsModel); err != nil {
		return Model.CmsID{}, err
	}
	return Model.CmsID{ID: cmsModel.ID}, nil
}

func ManagementService(cmsRepository Persistence.CrudRepository[Persistence.CmsModel]) IManagementService {
	return &managementService{cmsRepository: cmsRepository}
}
