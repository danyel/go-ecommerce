package productmanagement

import (
	Http "net/http"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	CommonHandler "github.com/danyel/ecommerce/internal/common/handler"
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

func (h *productManagementHandler) GetProducts(w Http.ResponseWriter, r *Http.Request) {
	products := h.s.GetProducts()
	CommonHandler.WriteResponse(Http.StatusOK, w, r, products)
}

func (h *productManagementHandler) DeleteProduct(w Http.ResponseWriter, r *Http.Request) {
	var productID Uuid.UUID
	var err error
	productIDToParse := Router.URLParam(r, "productID")
	if productID, err = Uuid.Parse(productIDToParse); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}

	if err = h.s.DeleteProduct(Types.NewID(productID)); err != nil {
		CommonHandler.StatusNotFound(w, r)
		return
	}
	CommonHandler.StatusNoContent(w, r)
}

func (h *productManagementHandler) UpdateProduct(w Http.ResponseWriter, r *Http.Request) {
	productID, err := CommonHandler.GetID(r, "productID")
	if err != nil {
		CommonHandler.StatusNotFound(w, r)
	}
	var updateProduct Product.UpdateProduct
	if err = CommonHandler.ValidateRequest(r, &updateProduct); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}
	if err = h.s.UpdateProduct(productID, updateProduct); err != nil {
		CommonHandler.StatusNotFound(w, r)
		return
	}
}

func (h *productManagementHandler) CreateProduct(w Http.ResponseWriter, r *Http.Request) {
	var createProduct Product.CreateProduct
	var productID Types.ID
	var err error

	if err = CommonHandler.ValidateRequest[Product.CreateProduct](r, &createProduct); err != nil {
		CommonHandler.StatusBadRequest(w, r)
	}

	if productID, err = h.s.CreateProduct(createProduct); err != nil {
		CommonHandler.StatusInternalServerError(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusCreated, w, r, productID)
}

func (h *productManagementHandler) GetProduct(w Http.ResponseWriter, r *Http.Request) {
	var productID Types.ID
	var err error
	var productModel Product.Product
	productID, err = CommonHandler.GetID(r, "productID")
	if err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}

	if productModel, err = h.s.GetProduct(productID); err != nil {
		CommonHandler.StatusNotFound(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusOK, w, r, productModel)
}

func NewHandler(DB *Database.DB) ProductManagementHandler {
	categoryService := Category.NewCategoryService(DB)
	cmsService := CMS.NewCmsService(DB)
	return &productManagementHandler{
		s: NewProductService(DB),
		m: Product.NewProductMapper(categoryService, cmsService),
	}
}
