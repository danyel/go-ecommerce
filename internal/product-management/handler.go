package product_management

import (
	Http "net/http"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	Router "github.com/go-chi/chi/v5"
	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductManagementHandler interface {
	GetProducts(w Http.ResponseWriter, _ *Http.Request)
	GetProduct(w Http.ResponseWriter, r *Http.Request)
	DeleteProduct(w Http.ResponseWriter, r *Http.Request)
	UpdateProduct(w Http.ResponseWriter, r *Http.Request)
	CreateProduct(w Http.ResponseWriter, r *Http.Request)
}

type productManagementHandler struct {
	s ProductService
	m Product.ProductMapper
}

func (h *productManagementHandler) GetProducts(w Http.ResponseWriter, _ *Http.Request) {
	products := h.s.GetProducts()
	commonHandler.WriteResponse(Http.StatusOK, w, products)
}

func (h *productManagementHandler) DeleteProduct(w Http.ResponseWriter, r *Http.Request) {
	var productId Uuid.UUID
	var err error
	productIdToParse := Router.URLParam(r, "productId")
	if productId, err = Uuid.Parse(productIdToParse); err != nil {
		commonHandler.StatusBadRequest(w)
		return
	}

	if err = h.s.DeleteProduct(Types.NewID(productId)); err != nil {
		commonHandler.StatusNotFound(w)
		return
	}
	commonHandler.StatusNoContent(w)
}
func (h *productManagementHandler) UpdateProduct(w Http.ResponseWriter, r *Http.Request) {
	productId, err := commonHandler.GetId(r, "productId")
	var updateProduct Product.UpdateProduct
	if err = commonHandler.ValidateRequest(r, &updateProduct); err != nil {
		commonHandler.StatusBadRequest(w)
		return
	}
	if err = h.s.UpdateProduct(productId, updateProduct); err != nil {
		commonHandler.StatusNotFound(w)
		return
	}
}
func (h *productManagementHandler) CreateProduct(w Http.ResponseWriter, r *Http.Request) {
	var createProduct Product.CreateProduct
	var productId Types.Id
	var err error

	if err = commonHandler.ValidateRequest[Product.CreateProduct](r, &createProduct); err != nil {
		commonHandler.StatusBadRequest(w)
	}

	if productId, err = h.s.CreateProduct(createProduct); err != nil {
		commonHandler.StatusInternalServerError(w)
		return
	}
	commonHandler.WriteResponse(Http.StatusCreated, w, productId)
}
func (h *productManagementHandler) GetProduct(w Http.ResponseWriter, r *Http.Request) {
	var productId Types.Id
	var err error
	var productModel Product.Product
	productId, err = commonHandler.GetId(r, "productId")
	if err != nil {
		commonHandler.StatusBadRequest(w)
		return
	}

	if productModel, err = h.s.GetProduct(productId); err != nil {
		commonHandler.StatusNotFound(w)
		return
	}
	commonHandler.WriteResponse(Http.StatusOK, w, productModel)
}

func NewHandler(DB *Database.DB) ProductManagementHandler {
	categoryService := Category.NewCategoryService(DB)
	cmsService := CMS.NewCmsService(DB)
	return &productManagementHandler{
		s: NewProductService(DB),
		m: Product.NewProductMapper(categoryService, cmsService),
	}
}
