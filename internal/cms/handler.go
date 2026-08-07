package cms

import (
	Http "net/http"

	CommonHandler "github.com/danyel/ecommerce/internal/common/handler"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsHandler interface {
	GetTranslation(w Http.ResponseWriter, r *Http.Request)
	GetTranslations(w Http.ResponseWriter, r *Http.Request)
}

type cmsHandler struct {
	s CmsService
}

func (h *cmsHandler) GetTranslation(w Http.ResponseWriter, r *Http.Request) {
	language := CommonHandler.GetPathParam(r, "language")
	code := CommonHandler.GetPathParam(r, "code")
	var translation Translation
	var err error

	if translation, err = h.s.GetTranslation(code, language); err != nil {
		CommonHandler.StatusNotFound(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, r, translation)
}

func (h *cmsHandler) GetTranslations(w Http.ResponseWriter, r *Http.Request) {
	language := CommonHandler.GetRequestParam(r, "language")

	CommonHandler.WriteResponse(Http.StatusOK, w, r, h.s.GetTranslations(language))
}

func NewHandler(db *Database.DB) CmsHandler {
	return &cmsHandler{
		NewCmsService(db),
	}
}
