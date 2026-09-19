package handler

import (
	Http "net/http"

	Model "github.com/danyel/ecommerce/internal/model"
	Service "github.com/danyel/ecommerce/internal/service"
	Types "github.com/danyel/ecommerce/internal/types"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryWebHandler interface {
	HandleCreateCategoryV1(response Http.ResponseWriter, request *Http.Request)
	HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request)
}

type categoryWebHandler struct {
	categoryService Service.ICategoryService
}

func (categoryWebHandler *categoryWebHandler) HandleCreateCategoryV1(response Http.ResponseWriter, request *Http.Request) {
	var createCategory Model.CreateCategory
	var categoryID Uuid.UUID
	var err error
	var details map[string]any
	if details, err = ValidateRequest(request, &createCategory); err != nil {
		BadRequest(response, request, InvalidRequestTitle, details)
		return
	}
	if categoryID, err = categoryWebHandler.categoryService.CreateCategory(createCategory); err != nil {
		InternalServerError(response, request, InternalServerErrorTitle, make(map[string]any))
		return
	}
	WriteResponse(Http.StatusCreated, response, request, Types.NewID(categoryID))
}

func (categoryWebHandler *categoryWebHandler) HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request) {
}

func NewCategoryWebHandler(categoryService Service.ICategoryService) CategoryWebHandler {
	handler := &categoryWebHandler{
		categoryService,
	}
	return handler
}
