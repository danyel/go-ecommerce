package management

import (
	Http "net/http"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	CommonHandler "github.com/danyel/ecommerce/internal/common/handler"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler interface {
	GetCategories(w Http.ResponseWriter, _ *Http.Request)
	CreateTranslations(w Http.ResponseWriter, _ *Http.Request)
}

type managementHandler struct {
	c   Category.CategoryService
	m   ManagementService
	cms CMS.CmsService
}

func (h *managementHandler) GetCategories(w Http.ResponseWriter, _ *Http.Request) {
	CommonHandler.WriteResponse(Http.StatusOK, w, h.c.GetCategories())
}

func (h *managementHandler) CreateTranslations(w Http.ResponseWriter, r *Http.Request) {
	var createCms CreateCms
	var err error
	var cmsId CmsId
	if err = CommonHandler.ValidateRequest[CreateCms](r, &createCms); err != nil {
		CommonHandler.StatusBadRequest(w)
		return
	}

	// we can not create a new translation for the same code and language!
	if _, err = h.cms.GetTranslation(createCms.Code, createCms.Language); err == nil {
		CommonHandler.StatusBadRequest(w)
	}

	if cmsId, err = h.m.CreateTranslation(createCms); err != nil {
		CommonHandler.StatusInternalServerError(w)
		return
	}
	CommonHandler.WriteResponse(Http.StatusCreated, w, cmsId)
}

func NewHandler(DB *Database.DB) ManagementHandler {
	return &managementHandler{
		c:   Category.NewCategoryService(DB),
		m:   NewManagementService(DB),
		cms: CMS.NewCmsService(DB),
	}
}
