package product

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductHandler interface {
	GetProducts(response Http.ResponseWriter, request *Http.Request)
	GetProduct(response Http.ResponseWriter, request *Http.Request)
}

type productHandler struct {
	productService ProductService
}

func (productHandler *productHandler) GetProducts(response Http.ResponseWriter, request *Http.Request) {
	WebHandler.WriteResponse(Http.StatusOK, response, request, productHandler.productService.GetProducts())
}

func (productHandler *productHandler) GetProduct(response Http.ResponseWriter, request *Http.Request) {
	var product Product
	var productID Types.ID
	var err error
	if productID, err = WebHandler.GetID(request, "productID"); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	if product, err = productHandler.productService.GetProduct(productID.ID); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, product)
}

func NewHandler(productService ProductService) ProductHandler {
	productHandler := &productHandler{productService}
	return productHandler
}
