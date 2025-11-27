package management

import (
	"encoding/json"
	"net/http"

	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler interface {
	GetCategories(w http.ResponseWriter, _ *http.Request)
	CreateTranslations(w http.ResponseWriter, _ *http.Request)
}

type managementHandler struct {
	categoryService   category.CategoryService
	managementService ManagementService
	handler           commonHandler.ResponseHandler
	cmsService        cms.CmsService
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
	if err = commonHandler.ValidateRequest[CreateCms](r, &createCms); err != nil {
		h.handler.StatusBadRequest(w)
		return
	}

	// we can not create a new translation for the same code and language!
	if _, err = h.cmsService.GetTranslation(createCms.Code, createCms.Language); err == nil {
		h.handler.StatusBadRequest(w)
	}

	if cmsId, err = h.managementService.CreateTranslation(createCms); err != nil {
		h.handler.StatusInternalServerError(w)
		return
	}
	h.handler.WriteResponse(http.StatusCreated, w, cmsId)
}

func NewHandler(DB *gorm.DB) ManagementHandler {
	return &managementHandler{
		categoryService:   category.NewCategoryService(DB),
		handler:           commonHandler.NewResponseHandler(),
		managementService: NewManagementService(DB),
		cmsService:        cms.NewCmsService(commonRepository.NewCrudRepository[cms.CmsModel](DB)),
	}
}
