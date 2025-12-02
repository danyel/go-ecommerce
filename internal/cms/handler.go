package cms

import (
	"net/http"

	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsHandler interface {
	GetTranslation(w http.ResponseWriter, r *http.Request)
	GetTranslations(w http.ResponseWriter, r *http.Request)
}

type cmsHandler struct {
	s CmsService
}

func (h *cmsHandler) GetTranslation(w http.ResponseWriter, r *http.Request) {
	language := commonHandler.GetPathParam(r, "language")
	code := commonHandler.GetPathParam(r, "code")
	var translation Translation
	var err error

	if translation, err = h.s.GetTranslation(code, language); err != nil {
		commonHandler.StatusNotFound(w)
		return
	}
	commonHandler.WriteResponse(http.StatusOK, w, translation)
}

func (h *cmsHandler) GetTranslations(w http.ResponseWriter, r *http.Request) {
	language := commonHandler.GetRequestParam(r, "language")

	commonHandler.WriteResponse(http.StatusOK, w, h.s.GetTranslations(language))
}

func NewHandler(db *gorm.DB) CmsHandler {
	return &cmsHandler{
		NewCmsService(db),
	}
}
