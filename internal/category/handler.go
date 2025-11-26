package category

import (
	"net/http"

	commonHandler "github.com/dnoulet/ecommerce/internal/common/handler"
	requestUtil "github.com/dnoulet/ecommerce/internal/util/request"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	CreateTranslations(_ http.ResponseWriter, _ *http.Request)
}

type categoryHandler struct {
	categoryService CategoryService
	Handler         commonHandler.ResponseHandler
}

func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var createCategory CreateCategory
	var categoryId CategoryId
	var err error
	if err = requestUtil.ValidateRequest(r, &createCategory); err != nil {
		h.Handler.StatusBadRequest(w)
		return
	}
	if categoryId, err = h.categoryService.CreateCategory(createCategory); err != nil {
		h.Handler.StatusInternalServerError(w)
		return
	}
	h.Handler.WriteResponse(http.StatusCreated, w, categoryId)
}

func (h *categoryHandler) CreateTranslations(_ http.ResponseWriter, _ *http.Request) {}

func NewHandler(DB *gorm.DB) CategoryHandler {
	handler := &categoryHandler{
		NewCategoryService(DB),
		commonHandler.NewResponseHandler(),
	}
	return handler
}
