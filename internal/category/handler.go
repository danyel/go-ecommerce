package category

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
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
	var categoryID CategoryID
	var err error
	if err = WebHandler.ValidateRequest(request, &createCategory); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if categoryID, err = categoryWebHandler.categoryService.CreateCategory(createCategory); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, categoryID)
}

func (categoryWebHandler *categoryWebHandler) HandleCreateTranslationsV1(response Http.ResponseWriter, request *Http.Request) {
}

func NewWebHandler(categoryService CategoryService) CategoryWebHandler {
	handler := &categoryWebHandler{
		categoryService,
	}
	return handler
}
