package cms

import (
	"net/http"

	"gorm.io/gorm"
)

type CmsHandler struct {
}

func (h *CmsHandler) GetTranslation(w http.ResponseWriter, r *http.Request) {}

func NewCmsHandler(DB *gorm.DB) CmsHandler {
	return CmsHandler{}
}
