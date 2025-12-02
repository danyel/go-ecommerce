package management

import (
	"net/http"

	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler interface {
	GetCategories(w http.ResponseWriter, _ *http.Request)
	CreateTranslations(w http.ResponseWriter, _ *http.Request)
}

type managementHandler struct {
	c   category.CategoryService
	m   ManagementService
	h   commonHandler.ResponseHandler
	cms cms.CmsService
}

func (h *managementHandler) GetCategories(w http.ResponseWriter, _ *http.Request) {
	h.h.WriteResponse(http.StatusOK, w, h.c.GetCategories())
}

func (h *managementHandler) CreateTranslations(w http.ResponseWriter, r *http.Request) {
	var createCms CreateCms
	var err error
	var cmsId CmsId
	if err = commonHandler.ValidateRequest[CreateCms](r, &createCms); err != nil {
		h.h.StatusBadRequest(w)
		return
	}

	// we can not create a new translation for the same code and language!
	if _, err = h.cms.GetTranslation(createCms.Code, createCms.Language); err == nil {
		h.h.StatusBadRequest(w)
	}

	if cmsId, err = h.m.CreateTranslation(createCms); err != nil {
		h.h.StatusInternalServerError(w)
		return
	}
	h.h.WriteResponse(http.StatusCreated, w, cmsId)
}

func NewHandler(DB *gorm.DB) ManagementHandler {
	return &managementHandler{
		c:   category.NewCategoryService(DB),
		h:   commonHandler.NewResponseHandler(),
		m:   NewManagementService(DB),
		cms: cms.NewCmsService(DB),
	}
}
