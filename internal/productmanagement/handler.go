package productmanagement

import (
	Http "net/http"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	Router "github.com/go-chi/chi/v5"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductManagementHandler interface {
	GetProducts(response Http.ResponseWriter, request *Http.Request)
	GetProduct(response Http.ResponseWriter, request *Http.Request)
	DeleteProduct(response Http.ResponseWriter, request *Http.Request)
	UpdateProduct(response Http.ResponseWriter, request *Http.Request)
	CreateProduct(response Http.ResponseWriter, request *Http.Request)
}

type productManagementHandler struct {
	productManagementService ProductManagementService
	productMapper            Product.ProductMapper
}

func (productManagementHandler *productManagementHandler) GetProducts(response Http.ResponseWriter, request *Http.Request) {
	products := productManagementHandler.productManagementService.GetProducts()
	WebHandler.WriteResponse(Http.StatusOK, response, request, products)
}

func (productManagementHandler *productManagementHandler) DeleteProduct(response Http.ResponseWriter, request *Http.Request) {
	var productID Uuid.UUID
	var err error
	productIDToParse := Router.URLParam(request, "productID")
	if productID, err = Uuid.Parse(productIDToParse); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	if err = productManagementHandler.productManagementService.DeleteProduct(Types.NewID(productID)); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.StatusNoContent(response, request)
}

func (productManagementHandler *productManagementHandler) UpdateProduct(response Http.ResponseWriter, request *Http.Request) {
	productID, err := WebHandler.GetID(request, "productID")
	if err != nil {
		WebHandler.StatusNotFound(response, request)
	}
	var updateProduct Product.UpdateProduct
	if err = WebHandler.ValidateRequest(request, &updateProduct); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if err = productManagementHandler.productManagementService.UpdateProduct(productID, updateProduct); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
}

func (productManagementHandler *productManagementHandler) CreateProduct(response Http.ResponseWriter, request *Http.Request) {
	var createProduct Product.CreateProduct
	var ID Types.ID
	var err error

	if err = WebHandler.ValidateRequest(request, &createProduct); err != nil {
		WebHandler.StatusBadRequest(response, request)
	}

	if ID, err = productManagementHandler.productManagementService.CreateProduct(createProduct); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, ID)
}

func (productManagementHandler *productManagementHandler) GetProduct(response Http.ResponseWriter, request *Http.Request) {
	var ID Types.ID
	var err error
	var product Product.Product
	ID, err = WebHandler.GetID(request, "productID")
	if err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	if product, err = productManagementHandler.productManagementService.GetProduct(ID); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, product)
}

func NewHandler(categoryService Category.CategoryService, cmsService CMS.CmsService, productManagementService ProductManagementService) ProductManagementHandler {
	return &productManagementHandler{
		productManagementService: productManagementService,
		productMapper:            Product.NewProductMapper(categoryService, cmsService),
	}
}
