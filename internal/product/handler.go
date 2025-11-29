package product

import (
	"encoding/json"
	"log"
	"net/http"

	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductHandler interface {
	GetProducts(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
}

type productApiHandler struct {
	productService ProductService
}
type productHtmlHandler struct {
	productService ProductService
}

func (h *productHtmlHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	err := Products(h.productService.GetProducts()).Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func (h *productHtmlHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := commonHandler.GetId(r, "productId")
	var product Product
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	if product, err = h.productService.GetProduct(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	err = ProductDetail(product).Render(r.Context(), w)
}

func (h *productApiHandler) GetProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.productService.GetProducts()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *productApiHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
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

func NewHtmlHandler(p ProductService) ProductHandler {
	h := &productHtmlHandler{p}
	return h
}
func NewApiHandler(p ProductService) ProductHandler {
	h := &productApiHandler{p}
	return h
}
