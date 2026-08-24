package management

import (
	Http "net/http"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementWebHandler interface {
	HandleGetCategoriesV1(response Http.ResponseWriter, request *Http.Request)
	HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request)
}

type managementWebHandler struct {
	categoryService   Category.CategoryService
	managementService ManagementService
	cmsService        CMS.CmsService
}

func (managementWebHandler *managementWebHandler) HandleGetCategoriesV1(response Http.ResponseWriter, request *Http.Request) {
	WebHandler.WriteResponse(Http.StatusOK, response, request, managementWebHandler.categoryService.GetCategories())
}

func (managementWebHandler *managementWebHandler) HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request) {
	var createCms CreateCms
	var err error
	var cmsID CmsID
	if err = WebHandler.ValidateRequest(request, &createCms); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	// we can not create a new translation for the same code and language!
	if _, err = managementWebHandler.cmsService.GetTranslation(createCms.Code, createCms.Language); err == nil {
		WebHandler.StatusBadRequest(response, request)
	}

	if cmsID, err = managementWebHandler.managementService.CreateTranslation(createCms); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, cmsID)
}

func NewWebHandler(categoryService Category.CategoryService, managementService ManagementService, cmsService CMS.CmsService) ManagementWebHandler {
	return &managementWebHandler{
		categoryService:   categoryService,
		managementService: managementService,
		cmsService:        cmsService,
	}
}
