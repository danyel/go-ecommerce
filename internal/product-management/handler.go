package product_management

import (
	"encoding/json"
	"net/http"

	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/internal/product"
	util "github.com/dnoulet/ecommerce/internal/util/request"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductManagementHandler interface {
	GetProducts(w http.ResponseWriter, _ *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
}

type productManagementHandler struct {
	productService ProductService
	productMapper  ProductMapper
}

func (h *productManagementHandler) GetProducts(w http.ResponseWriter, _ *http.Request) {
	products := h.productService.GetProducts()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.productMapper.MapProducts(products)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (h *productManagementHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var productId uuid.UUID
	var err error
	productIdToParse := chi.URLParam(r, "productId")
	if productId, err = uuid.Parse(productIdToParse); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.productService.DeleteProduct(productId); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *productManagementHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	parse, err := uuid.Parse(productId)
	var updateProduct UpdateProduct
	if err = util.ValidateRequest(r, &updateProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = h.productService.UpdateProduct(parse, updateProduct); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
func (h *productManagementHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var createProduct CreateProduct
	var productId ProductId
	var err error

	if err = util.ValidateRequest(r, &createProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if productId, err = h.productService.CreateProduct(createProduct); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(productId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (h *productManagementHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var productId uuid.UUID
	var err error
	var productModel product.Product
	productIdFromRequest := chi.URLParam(r, "productId")
	if productId, err = uuid.Parse(productIdFromRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if productModel, err = h.productService.GetProduct(productId); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.StatusText(200)
	if err = json.NewEncoder(w).Encode(h.productMapper.MapProduct(productModel)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewHandler(DB *gorm.DB) ProductManagementHandler {
	categoryService := category.NewCategoryService(DB)
	cmsService := cms.NewCmsService(commonRepository.NewCrudRepository[cms.CmsModel](DB))
	return &productManagementHandler{
		productService: NewProductService(DB),
		productMapper:  NewProductMapper(categoryService, cmsService),
	}
}
