package handler

import (
	Http "net/http"

	Domain "github.com/danyel/ecommerce/internal/domain"
	Model "github.com/danyel/ecommerce/internal/model"
	Category "github.com/danyel/ecommerce/internal/service"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementWebHandler interface {
	HandleGetCategoriesV1(response Http.ResponseWriter, request *Http.Request)
	HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request)
}

type managementWebHandler struct {
	categoryService   Category.ICategoryService
	managementService Category.IManagementService
	cmsService        Category.ICmsService
	maxPageSize       int
}

func (managementWebHandler *managementWebHandler) HandleGetCategoriesV1(response Http.ResponseWriter, request *Http.Request) {
	page, pageSize, err := ParsePagination(request, managementWebHandler.maxPageSize)
	if err != nil {
		BadRequest(response, request, InvalidRequestTitle, map[string]any{"pagination": err.Error()})
		return
	}
	categories, total := managementWebHandler.categoryService.GetCategoriesPage(page, pageSize)
	WriteResponse(Http.StatusOK, response, request, Model.Page[Domain.Category]{
		Items: categories, Page: page, PageSize: pageSize, Total: int(total),
		HasNextPage: page*pageSize < int(total),
	})
}

// HandleCreateTranslationsV1 still some messages to do here 😔
func (managementWebHandler *managementWebHandler) HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request) {
	var createCms Model.CreateCms
	var err error
	var cmsID Model.CmsID
	var details map[string]any
	if details, err = ValidateRequest(request, &createCms); err != nil {
		BadRequest(response, request, BadRequestTitle, details)
		return
	}

	// we can not create a new translation for the same code and language!
	if _, err = managementWebHandler.cmsService.GetTranslation(createCms.Code, createCms.Language); err == nil {
		BadRequest(response, request, BadRequestTitle, make(map[string]any))
	}

	if cmsID, err = managementWebHandler.managementService.CreateTranslation(createCms); err != nil {
		InternalServerError(response, request, InternalServerErrorTitle, make(map[string]any))
		return
	}
	WriteResponse(Http.StatusCreated, response, request, cmsID)
}

func NewManagementWebHandler(categoryService Category.ICategoryService, managementService Category.IManagementService, cmsService Category.ICmsService, maxPageSize ...int) ManagementWebHandler {
	maximum := 50
	if len(maxPageSize) > 0 && maxPageSize[0] > 0 {
		maximum = maxPageSize[0]
	}
	return &managementWebHandler{
		categoryService:   categoryService,
		managementService: managementService,
		cmsService:        cmsService,
		maxPageSize:       maximum,
	}
}
