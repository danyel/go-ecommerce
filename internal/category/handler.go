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
	var categoryID CategoryID
	var err error
	if err = CommonHandler.ValidateRequest(r, &createCategory); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}
	if categoryID, err = h.s.CreateCategory(createCategory); err != nil {
		CommonHandler.StatusInternalServerError(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusCreated, w, r, categoryID)
}

func (h *categoryHandler) CreateTranslations(_ Http.ResponseWriter, _ *Http.Request) {}

func NewHandler(DB *Database.DB) CategoryHandler {
	handler := &categoryHandler{
		NewCategoryService(DB),
	}
	return handler
}
