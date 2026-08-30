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
type ProductManagementWebHandler interface {
	HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetProductV1(response Http.ResponseWriter, request *Http.Request)
	HandleDeleteProductV1(response Http.ResponseWriter, request *Http.Request)
	HandleUpdateProductV1(response Http.ResponseWriter, request *Http.Request)
	HandleCreateProductV1(response Http.ResponseWriter, request *Http.Request)
}

type productManagementWebHandler struct {
	productManagementService ProductManagementService
	productMapper            Product.ProductMapper
}

func (productManagementWebHandler *productManagementWebHandler) HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request) {
	products := productManagementWebHandler.productManagementService.GetProducts()
	WebHandler.WriteResponse(Http.StatusOK, response, request, products)
}

func (productManagementWebHandler *productManagementWebHandler) HandleDeleteProductV1(response Http.ResponseWriter, request *Http.Request) {
	var productID Uuid.UUID
	var err error
	productIDToParse := Router.URLParam(request, "productID")
	if productID, err = Uuid.Parse(productIDToParse); err != nil {
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, make(map[string]any))
		return
	}

	if err = productManagementWebHandler.productManagementService.DeleteProduct(Types.NewID(productID)); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.StatusNoContent(response, request)
}

func (productManagementWebHandler *productManagementWebHandler) HandleUpdateProductV1(response Http.ResponseWriter, request *Http.Request) {
	productID, err := WebHandler.GetID(request)
	if err != nil {
		WebHandler.StatusNotFound(response, request)
	}
	var updateProduct Product.UpdateProductDTO
	var details map[string]any
	if details, err = WebHandler.ValidateRequest(request, &updateProduct); err != nil {
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, details)
		return
	}
	if err = productManagementWebHandler.productManagementService.UpdateProduct(productID, updateProduct); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
}

func (productManagementWebHandler *productManagementWebHandler) HandleCreateProductV1(response Http.ResponseWriter, request *Http.Request) {
	var createProduct Product.CreateProductDTO
	var ID Types.ID
	var err error
	var details map[string]any

	if details, err = WebHandler.ValidateRequest(request, &createProduct); err != nil {
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, details)
	}

	if ID, err = productManagementWebHandler.productManagementService.CreateProduct(createProduct); err != nil {
		WebHandler.InternalServerError(response, request, WebHandler.InternalServerErrorTitle, make(map[string]any))
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, ID)
}

func (productManagementWebHandler *productManagementWebHandler) HandleGetProductV1(response Http.ResponseWriter, request *Http.Request) {
	var ID Types.ID
	var err error
	var product Product.Product
	ID, err = WebHandler.GetID(request)
	if err != nil {
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, make(map[string]any))
		return
	}

	if product, err = productManagementWebHandler.productManagementService.GetProduct(ID); err != nil {
		WebHandler.StatusNotFound(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, dto(product))
}

func dto(product Product.Product) Product.ProductDTO {
	return Product.ProductDTO{
		ID:          product.ID,
		Brand:       product.Brand,
		Name:        product.Name,
		Description: product.Description,
		Code:        product.Code,
		Price:       product.Price.DTO(),
		Category:    product.Category,
		ImageURL:    product.ImageURL,
		Stock:       product.Stock,
	}
}

func NewWebHandler(categoryService Category.CategoryService, cmsService CMS.CmsService, productManagementService ProductManagementService) ProductManagementWebHandler {
	return &productManagementWebHandler{
		productManagementService: productManagementService,
		productMapper:            Product.NewProductMapper(categoryService, cmsService),
	}
}
