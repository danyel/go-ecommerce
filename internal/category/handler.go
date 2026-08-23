package category

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryHandler interface {
	CreateCategory(response Http.ResponseWriter, request *Http.Request)
	CreateTranslations(response Http.ResponseWriter, request *Http.Request)
}

type categoryHandler struct {
	categoryService CategoryService
}

func (categoryHandler *categoryHandler) CreateCategory(response Http.ResponseWriter, request *Http.Request) {
	var createCategory CreateCategory
	var categoryID CategoryID
	var err error
	if err = WebHandler.ValidateRequest(request, &createCategory); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if categoryID, err = categoryHandler.categoryService.CreateCategory(createCategory); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, categoryID)
}

func (categoryHandler *categoryHandler) CreateTranslations(response Http.ResponseWriter, request *Http.Request) {
}

func NewHandler(categoryService CategoryService) CategoryHandler {
	handler := &categoryHandler{
		categoryService,
	}
	return handler
}
