package product

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductHandler struct {
	GetProducts func(w http.ResponseWriter, _ *http.Request)
	GetProduct  func(w http.ResponseWriter, r *http.Request)
}

func NewHandler(DB *gorm.DB) ProductHandler {
	productService := NewProductService(DB)
	return ProductHandler{
		GetProducts: func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(productService.GetProducts()); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
		GetProduct: func(w http.ResponseWriter, r *http.Request) {
			var product Product
			var productId uuid.UUID
			var err error
			productIdToParse := chi.URLParam(r, "productId")
			if productId, err = uuid.Parse(productIdToParse); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if product, err = productService.GetProduct(productId); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			http.StatusText(200)

			if err = json.NewEncoder(w).Encode(product); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
	}
}
