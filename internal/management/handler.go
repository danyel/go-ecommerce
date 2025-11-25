package management

import (
	"encoding/json"
	"net/http"

	"github.com/dnoulet/ecommerce/internal/category"
	commonHandler "github.com/dnoulet/ecommerce/internal/common/handler"
	requestutil "github.com/dnoulet/ecommerce/internal/util/request"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler interface {
	GetCategories(w http.ResponseWriter, _ *http.Request)
	CreateTranslations(w http.ResponseWriter, _ *http.Request)
}

type managementHandler struct {
	categoryService category.CategoryService
	cmsService      CmsService
	handler         commonHandler.ResponseHandler
}

func (h *managementHandler) GetCategories(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.handler.StatusOK(w)
	if err := json.NewEncoder(w).Encode(h.categoryService.GetCategories()); err != nil {
		h.handler.StatusInternalServerError(w)
	}
}

func (h *managementHandler) CreateTranslations(w http.ResponseWriter, r *http.Request) {
	var createCms CreateCms
	var err error
	var cmsId CmsId
	if err = requestutil.ValidateRequest[CreateCms](r, &createCms); err != nil {
		h.handler.StatusBadRequest(w)
		return
	}
	if cmsId, err = h.cmsService.CreateTranslation(createCms); err != nil {
		h.handler.StatusInternalServerError(w)
		return
	}
	h.handler.WriteResponse(w, cmsId)
}

func NewHandler(DB *gorm.DB) ManagementHandler {
	return &managementHandler{
		categoryService: category.NewCategoryService(DB),
		handler:         commonHandler.NewResponseHandler(),
		cmsService:      NewCmsService(DB),
	}
}
