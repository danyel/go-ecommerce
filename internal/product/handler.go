package product

import (
	Http "net/http"

	CommonHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductHandler interface {
	GetProducts(w Http.ResponseWriter, r *Http.Request)
	GetProduct(w Http.ResponseWriter, r *Http.Request)
}

type productAPIHandler struct {
	p ProductService
}

func (h *productAPIHandler) GetProducts(w Http.ResponseWriter, r *Http.Request) {
	CommonHandler.WriteResponse(Http.StatusOK, w, r, h.p.GetProducts())
}

func (h *productAPIHandler) GetProduct(w Http.ResponseWriter, r *Http.Request) {
	var product Product
	var productID Types.ID
	var err error
	if productID, err = CommonHandler.GetID(r, "productID"); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}

	if product, err = h.p.GetProduct(productID.ID); err != nil {
		CommonHandler.StatusNotFound(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, r, product)
}

func NewAPIHandler(p ProductService) ProductHandler {
	h := &productAPIHandler{p}
	return h
}
