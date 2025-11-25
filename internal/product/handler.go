package product

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductHandler interface {
	GetProducts(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	productService ProductService
}

func (h *productHandler) GetProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.productService.GetProducts()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *productHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	var productId uuid.UUID
	var err error
	productIdToParse := chi.URLParam(r, "productId")
	if productId, err = uuid.Parse(productIdToParse); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if product, err = h.productService.GetProduct(productId); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.StatusText(200)

	if err = json.NewEncoder(w).Encode(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewHandler(DB *gorm.DB) ProductHandler {
	h := &productHandler{NewProductService(DB)}
	return h
}
