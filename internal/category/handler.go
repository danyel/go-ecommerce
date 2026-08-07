package category

import (
	Http "net/http"

	CommonHandler "github.com/danyel/ecommerce/internal/common/handler"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryHandler interface {
	CreateCategory(w Http.ResponseWriter, r *Http.Request)
	CreateTranslations(_ Http.ResponseWriter, _ *Http.Request)
}

type categoryHandler struct {
	s CategoryService
}

func (h *categoryHandler) CreateCategory(w Http.ResponseWriter, r *Http.Request) {
	var createCategory CreateCategory
	var categoryId CategoryId
	var err error
	if err = CommonHandler.ValidateRequest[CreateCategory](r, &createCategory); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}
	if categoryId, err = h.s.CreateCategory(createCategory); err != nil {
		CommonHandler.StatusInternalServerError(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusCreated, w, r, categoryId)
}

func (h *categoryHandler) CreateTranslations(_ Http.ResponseWriter, _ *Http.Request) {}

func NewHandler(DB *Database.DB) CategoryHandler {
	handler := &categoryHandler{
		NewCategoryService(DB),
	}
	return handler
}
