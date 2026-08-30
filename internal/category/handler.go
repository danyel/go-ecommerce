package category

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryWebHandler interface {
	HandleCreateCategoryV1(response Http.ResponseWriter, request *Http.Request)
	HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request)
}

type categoryWebHandler struct {
	categoryService CategoryService
}

func (categoryWebHandler *categoryWebHandler) HandleCreateCategoryV1(response Http.ResponseWriter, request *Http.Request) {
	var createCategory CreateCategory
	var categoryID Uuid.UUID
	var err error
	var details map[string]any
	if details, err = WebHandler.ValidateRequest(request, &createCategory); err != nil {
		WebHandler.BadRequest(response, request, WebHandler.InvalidRequestTitle, details)
		return
	}
	if categoryID, err = categoryWebHandler.categoryService.CreateCategory(createCategory); err != nil {
		WebHandler.InternalServerError(response, request, WebHandler.InternalServerErrorTitle, make(map[string]any))
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, Types.NewID(categoryID))
}

func (categoryWebHandler *categoryWebHandler) HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request) {
}

func NewWebHandler(categoryService CategoryService) CategoryWebHandler {
	handler := &categoryWebHandler{
		categoryService,
	}
	return handler
}
