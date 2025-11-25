package cms

import (
	"net/http"

	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsHandler interface {
	GetTranslation(w http.ResponseWriter, r *http.Request)
}

type cmsHandler struct{}

func (h *cmsHandler) GetTranslation(_ http.ResponseWriter, _ *http.Request) {}

func NewHandler(_ *gorm.DB) CmsHandler {
	return &cmsHandler{}
}
