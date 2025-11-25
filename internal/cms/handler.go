package cms

import (
	"net/http"

	commonHandler "github.com/dnoulet/ecommerce/internal/common/handler"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsHandler interface {
	GetTranslation(w http.ResponseWriter, r *http.Request)
	GetTranslations(w http.ResponseWriter, r *http.Request)
}

type cmsHandler struct {
	cmsService CmsService
	Handler    commonHandler.ResponseHandler
}

func (h *cmsHandler) GetTranslation(w http.ResponseWriter, r *http.Request) {
	language := chi.URLParam(r, "language")
	var cmsId uuid.UUID
	var translation Translation
	var err error

	if cmsId, err = uuid.Parse(chi.URLParam(r, "id")); err != nil {
		h.Handler.StatusBadRequest(w)
		return
	}
	if translation, err = h.cmsService.GetTranslation(cmsId, language); err != nil {
		h.Handler.StatusNotFound(w)
		return
	}
	h.Handler.WriteResponse(w, translation)
}

func (h *cmsHandler) GetTranslations(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")

	h.Handler.WriteResponse(w, h.cmsService.GetTranslations(language))
}

func NewHandler(db *gorm.DB) CmsHandler {
	return &cmsHandler{
		NewCmsService(db),
		commonHandler.NewResponseHandler(),
	}
}
