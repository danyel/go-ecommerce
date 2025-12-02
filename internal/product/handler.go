package product

import (
	"net/http"

	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductHandler interface {
	GetProducts(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
}

type productApiHandler struct {
	p ProductService
	h commonHandler.ResponseHandler
}

func (h *productApiHandler) GetProducts(w http.ResponseWriter, _ *http.Request) {
	h.h.WriteResponse(http.StatusOK, w, h.p.GetProducts())
}

func (h *productApiHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	var productId uuid.UUID
	var err error
	if productId, err = commonHandler.GetId(r, "productId"); err != nil {
		h.h.StatusBadRequest(w)
		return
	}

	if product, err = h.p.GetProduct(productId); err != nil {
		h.h.StatusNotFound(w)
		return
	}
	h.h.WriteResponse(http.StatusOK, w, product)
}

func NewApiHandler(p ProductService) ProductHandler {
	h := &productApiHandler{p, commonHandler.NewResponseHandler()}
	return h
}
