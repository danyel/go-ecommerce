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
	cms cms.CmsService
}

func (h *managementHandler) GetCategories(w http.ResponseWriter, _ *http.Request) {
	commonHandler.WriteResponse(http.StatusOK, w, h.c.GetCategories())
}

func (h *managementHandler) CreateTranslations(w http.ResponseWriter, r *http.Request) {
	var createCms CreateCms
	var err error
	var cmsId CmsId
	if err = commonHandler.ValidateRequest[CreateCms](r, &createCms); err != nil {
		commonHandler.StatusBadRequest(w)
		return
	}

	// we can not create a new translation for the same code and language!
	if _, err = h.cms.GetTranslation(createCms.Code, createCms.Language); err == nil {
		commonHandler.StatusBadRequest(w)
	}

	if cmsId, err = h.m.CreateTranslation(createCms); err != nil {
		commonHandler.StatusInternalServerError(w)
		return
	}
	commonHandler.WriteResponse(http.StatusCreated, w, cmsId)
}

func NewHandler(DB *gorm.DB) ManagementHandler {
	return &managementHandler{
		c:   category.NewCategoryService(DB),
		m:   NewManagementService(DB),
		cms: cms.NewCmsService(DB),
	}
}
