package cms

import (
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsWebHandler interface {
	HandleGetTranslationV1(response Http.ResponseWriter, request *Http.Request)
	HandleV1(response Http.ResponseWriter, request *Http.Request)
}

type cmsWebHandler struct {
	cmsService CmsService
}

func (cmsWebHandler *cmsWebHandler) HandleGetTranslationV1(response Http.ResponseWriter, request *Http.Request) {
	language := WebHandler.GetPathParam(request, "language")
	code := WebHandler.GetPathParam(request, "code")
	var translation Translation
	var err error

	if translation, err = cmsWebHandler.cmsService.GetTranslation(code, language); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, translation)
}

func (cmsWebHandler *cmsWebHandler) HandleV1(response Http.ResponseWriter, request *Http.Request) {
	language := WebHandler.GetPathParam(request, "language")
	Logger.Log.Debug("Searching translations for language: %s", language)
	WebHandler.WriteResponse(Http.StatusOK, response, request, cmsWebHandler.cmsService.GetTranslations(language))
}

func NewWebHandler(cmsService CmsService) CmsWebHandler {
	return &cmsWebHandler{
		cmsService,
	}
}
