package cms

import (
	"net/http"

	commonHandler "github.com/dnoulet/ecommerce/internal/common/handler"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/go-chi/chi/v5"
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
	code := chi.URLParam(r, "code")
	var translation Translation
	var err error

	if translation, err = h.cmsService.GetTranslation(code, language); err != nil {
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
		NewCmsService(commonRepository.NewCrudRepository[CmsModel](db)),
		commonHandler.NewResponseHandler(),
	}
}
