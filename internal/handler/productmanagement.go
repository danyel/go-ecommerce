package handler

import (
	Http "net/http"

	Domain "github.com/danyel/ecommerce/internal/domain"
	Mapper "github.com/danyel/ecommerce/internal/mapper"
	Model "github.com/danyel/ecommerce/internal/model"
	Service "github.com/danyel/ecommerce/internal/service"
	Types "github.com/danyel/ecommerce/internal/types"
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
	productManagementService Service.IProductManagementService
	productMapper            Mapper.ProductMapper
	categoryMapper           Mapper.ICategoryMapper
	maxPageSize              int
}

func (productManagementWebHandler *productManagementWebHandler) HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request) {
	page, pageSize, err := ParsePagination(request, productManagementWebHandler.maxPageSize)
	if err != nil {
		BadRequest(response, request, InvalidRequestTitle, map[string]any{"pagination": err.Error()})
		return
	}
	products, total := productManagementWebHandler.productManagementService.GetProductsPage(page, pageSize)
	WriteResponse(Http.StatusOK, response, request, Model.Page[Domain.Product]{
		Items: products, Page: page, PageSize: pageSize, Total: int(total),
		HasNextPage: page*pageSize < int(total),
	})
}

func (productManagementWebHandler *productManagementWebHandler) HandleDeleteProductV1(response Http.ResponseWriter, request *Http.Request) {
	var productID Uuid.UUID
	var err error
	productIDToParse := Router.URLParam(request, "ID")
	if productID, err = Uuid.Parse(productIDToParse); err != nil {
		BadRequest(response, request, BadRequestTitle, make(map[string]any))
		return
	}

	if err = productManagementWebHandler.productManagementService.DeleteProduct(Types.NewID(productID)); err != nil {
		StatusNotFound(response, request)
		return
	}
	StatusNoContent(response, request)
}

func (productManagementWebHandler *productManagementWebHandler) HandleUpdateProductV1(response Http.ResponseWriter, request *Http.Request) {
	productID, err := GetID(request)
	if err != nil {
		StatusNotFound(response, request)
		return
	}
	var updateProduct Model.UpdateProductDTO
	var details map[string]any
	if details, err = ValidateRequest(request, &updateProduct); err != nil {
		BadRequest(response, request, BadRequestTitle, details)
		return
	}
	if err = productManagementWebHandler.productManagementService.UpdateProduct(productID, updateProduct); err != nil {
		StatusNotFound(response, request)
		return
	}
}

func (productManagementWebHandler *productManagementWebHandler) HandleCreateProductV1(response Http.ResponseWriter, request *Http.Request) {
	var createProduct Model.CreateProductDTO
	var ID Types.ID
	var err error
	var details map[string]any

	if details, err = ValidateRequest(request, &createProduct); err != nil {
		BadRequest(response, request, BadRequestTitle, details)
		return
	}

	if ID, err = productManagementWebHandler.productManagementService.CreateProduct(createProduct); err != nil {
		InternalServerError(response, request, InternalServerErrorTitle, make(map[string]any))
		return
	}
	WriteResponse(Http.StatusCreated, response, request, ID)
}

func (productManagementWebHandler *productManagementWebHandler) HandleGetProductV1(response Http.ResponseWriter, request *Http.Request) {
	var ID Types.ID
	var err error
	var product Domain.Product
	ID, err = GetID(request)
	if err != nil {
		BadRequest(response, request, BadRequestTitle, make(map[string]any))
		return
	}

	if product, err = productManagementWebHandler.productManagementService.GetProduct(ID); err != nil {
		StatusNotFound(response, request)
		return
	}
	WriteResponse(Http.StatusOK, response, request, productManagementWebHandler.productDTO(product))
}

func (productManagementWebHandler *productManagementWebHandler) productDTO(product Domain.Product) Model.ProductDTO {
	return Model.ProductDTO{
		ID:          product.ID,
		Brand:       product.Brand,
		Name:        product.Name,
		Description: product.Description,
		Code:        product.Code,
		Price:       product.Price.DTO(),
		Category:    productManagementWebHandler.categoryMapper.Map(product.Category),
		ImageURL:    product.ImageURL,
		Stock:       product.Stock,
	}
}

func NewProductManagementWebHandler(categoryService Service.ICategoryService, cmsService Service.ICmsService, productManagementService Service.IProductManagementService, productMapper Mapper.ProductMapper, categoryMapper Mapper.ICategoryMapper, maxPageSize ...int) ProductManagementWebHandler {
	maximum := 50
	if len(maxPageSize) > 0 && maxPageSize[0] > 0 {
		maximum = maxPageSize[0]
	}
	return &productManagementWebHandler{
		productManagementService: productManagementService,
		productMapper:            productMapper,
		categoryMapper:           categoryMapper,
		maxPageSize:              maximum,
	}
}
