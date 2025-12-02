package category

import (
	"net/http"

	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	CreateTranslations(_ http.ResponseWriter, _ *http.Request)
}

type categoryHandler struct {
	s CategoryService
}

func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var createCategory CreateCategory
	var categoryId CategoryId
	var err error
	if err = commonHandler.ValidateRequest[CreateCategory](r, &createCategory); err != nil {
		commonHandler.StatusBadRequest(w)
		return
	}
	if categoryId, err = h.s.CreateCategory(createCategory); err != nil {
		commonHandler.StatusInternalServerError(w)
		return
	}
	commonHandler.WriteResponse(http.StatusCreated, w, categoryId)
}

func (h *categoryHandler) CreateTranslations(_ http.ResponseWriter, _ *http.Request) {}

func NewHandler(DB *gorm.DB) CategoryHandler {
	handler := &categoryHandler{
		NewCategoryService(DB),
	}
	return handler
}
