package cms

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsHandler interface {
	GetTranslation(response Http.ResponseWriter, request *Http.Request)
	GetTranslations(response Http.ResponseWriter, request *Http.Request)
}

type cmsHandler struct {
	cmsService CmsService
}

func (cmsHandler *cmsHandler) GetTranslation(response Http.ResponseWriter, request *Http.Request) {
	language := WebHandler.GetPathParam(request, "language")
	code := WebHandler.GetPathParam(request, "code")
	var translation Translation
	var err error

	if translation, err = cmsHandler.cmsService.GetTranslation(code, language); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, translation)
}

func (cmsHandler *cmsHandler) GetTranslations(response Http.ResponseWriter, request *Http.Request) {
	language := WebHandler.GetRequestParam(request, "language")

	WebHandler.WriteResponse(Http.StatusOK, response, request, cmsHandler.cmsService.GetTranslations(language))
}

func NewHandler(cmsService CmsService) CmsHandler {
	return &cmsHandler{
		cmsService,
	}
}
