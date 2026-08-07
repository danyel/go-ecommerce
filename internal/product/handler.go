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

type productApiHandler struct {
	p ProductService
}

func (h *productApiHandler) GetProducts(w Http.ResponseWriter, r *Http.Request) {
	CommonHandler.WriteResponse(Http.StatusOK, w, r, h.p.GetProducts())
}

func (h *productApiHandler) GetProduct(w Http.ResponseWriter, r *Http.Request) {
	var product Product
	var productId Types.Id
	var err error
	if productId, err = CommonHandler.GetId(r, "productId"); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}

	if product, err = h.p.GetProduct(productId.ID); err != nil {
		CommonHandler.StatusNotFound(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, r, product)
}

func NewApiHandler(p ProductService) ProductHandler {
	h := &productApiHandler{p}
	return h
}
