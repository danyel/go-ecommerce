package product_management

import (
	"net/http"

	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"github.com/danyel/ecommerce/internal/product"
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
	s ProductService
	m product.ProductMapper
	h commonHandler.ResponseHandler
}

func (h *productManagementHandler) GetProducts(w http.ResponseWriter, _ *http.Request) {
	products := h.s.GetProducts()
	h.h.WriteResponse(http.StatusOK, w, products)
}

func (h *productManagementHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var productId uuid.UUID
	var err error
	productIdToParse := chi.URLParam(r, "productId")
	if productId, err = uuid.Parse(productIdToParse); err != nil {
		h.h.StatusBadRequest(w)
		return
	}

	if err = h.s.DeleteProduct(productId); err != nil {
		h.h.StatusNotFound(w)
		return
	}
	h.h.StatusNoContent(w)
}
func (h *productManagementHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := commonHandler.GetId(r, "productId")
	var updateProduct UpdateProduct
	if err = commonHandler.ValidateRequest(r, &updateProduct); err != nil {
		h.h.StatusBadRequest(w)
		return
	}
	if err = h.s.UpdateProduct(productId, updateProduct); err != nil {
		h.h.StatusNotFound(w)
		return
	}
}
func (h *productManagementHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var createProduct CreateProduct
	var productId ProductId
	var err error

	if err = commonHandler.ValidateRequest[CreateProduct](r, &createProduct); err != nil {
		h.h.StatusBadRequest(w)
	}

	if productId, err = h.s.CreateProduct(createProduct); err != nil {
		h.h.StatusInternalServerError(w)
		return
	}
	h.h.WriteResponse(http.StatusCreated, w, productId)
}
func (h *productManagementHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var productId uuid.UUID
	var err error
	var productModel product.Product
	productId, err = commonHandler.GetId(r, "productId")
	if err != nil {
		h.h.StatusBadRequest(w)
		return
	}

	if productModel, err = h.s.GetProduct(productId); err != nil {
		h.h.StatusNotFound(w)
		return
	}
	h.h.WriteResponse(http.StatusOK, w, productModel)
}

func NewHandler(DB *gorm.DB) ProductManagementHandler {
	categoryService := category.NewCategoryService(DB)
	cmsService := cms.NewCmsService(DB)
	return &productManagementHandler{
		s: NewProductService(DB),
		m: product.NewProductMapper(categoryService, cmsService),
		h: commonHandler.NewResponseHandler(),
	}
}
