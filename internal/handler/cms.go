package handler

import (
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Model "github.com/danyel/ecommerce/internal/model"
	Service "github.com/danyel/ecommerce/internal/service"
)

//goland:noinspection GoNameStartsWithPackageName
type CmsWebHandler interface {
	HandleGetTranslationV1(response Http.ResponseWriter, request *Http.Request)
	HandleV1(response Http.ResponseWriter, request *Http.Request)
}

type cmsWebHandler struct {
	cmsService Service.ICmsService
}

func (cmsWebHandler *cmsWebHandler) HandleGetTranslationV1(response Http.ResponseWriter, request *Http.Request) {
	language := GetPathParam(request, "language")
	code := GetPathParam(request, "code")
	var translation Model.Translation
	var err error

	if translation, err = cmsWebHandler.cmsService.GetTranslation(code, language); err != nil {
		StatusNotFound(response, request)
		return
	}
	WriteResponse(Http.StatusOK, response, request, translation)
}

func (cmsWebHandler *cmsWebHandler) HandleV1(response Http.ResponseWriter, request *Http.Request) {
	language := GetPathParam(request, "language")
	Logger.Log.Debug("Searching translations for language: %s", language)
	WriteResponse(Http.StatusOK, response, request, cmsWebHandler.cmsService.GetTranslations(language))
}

func NewCmsWebHandler(cmsService Service.ICmsService) CmsWebHandler {
	return &cmsWebHandler{
		cmsService,
	}
}
