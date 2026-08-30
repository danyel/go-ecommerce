package product

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductWebHandler interface {
	HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetProductV1(response Http.ResponseWriter, request *Http.Request)
}

type productWebHandler struct {
	productService ProductService
}

func (productWebHandler *productWebHandler) HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request) {
	WebHandler.WriteResponse(Http.StatusOK, response, request, productWebHandler.productService.GetProducts())
}

func (productWebHandler *productWebHandler) HandleGetProductV1(response Http.ResponseWriter, request *Http.Request) {
	var product Product
	var ID Types.ID
	var err error
	if ID, err = WebHandler.GetID(request); err != nil {
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, make(map[string]any))
		return
	}

	if product, err = productWebHandler.productService.GetProduct(ID.ID); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, product)
}

func NewWebHandler(productService ProductService) ProductWebHandler {
	productWebHandler := &productWebHandler{productService}
	return productWebHandler
}
