package management

import (
	Http "net/http"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler interface {
	GetCategories(response Http.ResponseWriter, request *Http.Request)
	CreateTranslations(response Http.ResponseWriter, request *Http.Request)
}

type managementHandler struct {
	categoryService   Category.CategoryService
	managementService ManagementService
	cmsService        CMS.CmsService
}

func (managementHandler *managementHandler) GetCategories(response Http.ResponseWriter, request *Http.Request) {
	WebHandler.WriteResponse(Http.StatusOK, response, request, managementHandler.categoryService.GetCategories())
}

func (managementHandler *managementHandler) CreateTranslations(response Http.ResponseWriter, request *Http.Request) {
	var createCms CreateCms
	var err error
	var cmsID CmsID
	if err = WebHandler.ValidateRequest(request, &createCms); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	// we can not create a new translation for the same code and language!
	if _, err = managementHandler.cmsService.GetTranslation(createCms.Code, createCms.Language); err == nil {
		WebHandler.StatusBadRequest(response, request)
	}

	if cmsID, err = managementHandler.managementService.CreateTranslation(createCms); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, cmsID)
}

func NewHandler(categoryService Category.CategoryService, managementService ManagementService, cmsService CMS.CmsService) ManagementHandler {
	return &managementHandler{
		categoryService:   categoryService,
		managementService: managementService,
		cmsService:        cmsService,
	}
}
