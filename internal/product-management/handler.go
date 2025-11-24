package product_management

import (
	"encoding/json"
	"net/http"

	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/product"
	util "github.com/dnoulet/ecommerce/internal/util/request"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductManagementHandler struct {
	GetProducts   func(w http.ResponseWriter, _ *http.Request)
	GetProduct    func(w http.ResponseWriter, r *http.Request)
	DeleteProduct func(w http.ResponseWriter, r *http.Request)
	UpdateProduct func(w http.ResponseWriter, r *http.Request)
	CreateProduct func(w http.ResponseWriter, r *http.Request)
}

func NewHandler(DB *gorm.DB) ProductManagementHandler {
	categoryService := category.NewCategoryService(DB)
	productService := NewProductService(DB)
	productMapper := NewProductMapper(&categoryService)
	return ProductManagementHandler{
		GetProducts: func(w http.ResponseWriter, _ *http.Request) {
			products := productService.GetProducts()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(productMapper.MapProducts(products)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
		DeleteProduct: func(w http.ResponseWriter, r *http.Request) {
			var productId uuid.UUID
			var err error
			productIdToParse := chi.URLParam(r, "productId")
			if productId, err = uuid.Parse(productIdToParse); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = productService.DeleteProduct(productId); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
		UpdateProduct: func(w http.ResponseWriter, r *http.Request) {
			productId := chi.URLParam(r, "productId")
			parse, err := uuid.Parse(productId)
			var updateProduct UpdateProduct
			if err = util.ValidateRequest(r, &updateProduct); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err = productService.UpdateProduct(parse, updateProduct); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		},
		CreateProduct: func(w http.ResponseWriter, r *http.Request) {
			var createProduct CreateProduct
			var productId ProductId
			var err error

			if err = util.ValidateRequest(r, &createProduct); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			if productId, err = productService.CreateProduct(createProduct); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err = json.NewEncoder(w).Encode(productId); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
		GetProduct: func(w http.ResponseWriter, r *http.Request) {
			var productId uuid.UUID
			var err error
			var productModel product.Product
			productIdFromRequest := chi.URLParam(r, "productId")
			if productId, err = uuid.Parse(productIdFromRequest); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if productModel, err = productService.GetProduct(productId); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			http.StatusText(200)
			if err = json.NewEncoder(w).Encode(productMapper.MapProduct(productModel)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
	}
}
