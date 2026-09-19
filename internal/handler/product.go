package handler

import (
	Http "net/http"

	Domain "github.com/danyel/ecommerce/internal/domain"
	Mapper "github.com/danyel/ecommerce/internal/mapper"
	Model "github.com/danyel/ecommerce/internal/model"
	Service "github.com/danyel/ecommerce/internal/service"
	Types "github.com/danyel/ecommerce/internal/types"
)

//goland:noinspection GoNameStartsWithPackageName
type IProductWebHandler interface {
	HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetProductV1(response Http.ResponseWriter, request *Http.Request)
}

type productWebHandler struct {
	productService Service.IProductService
	productMapper  Mapper.ProductMapper
	categoryMapper Mapper.ICategoryMapper
}

func (productWebHandler *productWebHandler) HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request) {
	products := productWebHandler.productService.FindAll()
	productModels := make([]Model.ProductDTO, len(products))
	for i, product := range products {
		productModels[i] = Model.ProductDTO{
			ID:          product.ID,
			Brand:       product.Brand,
			Name:        product.Name,
			Description: product.Description,
			Code:        product.Code,
			Price:       product.Price.DTO(),
			Category:    productWebHandler.categoryMapper.Map(product.Category),
			ImageURL:    product.ImageURL,
			Stock:       product.Stock,
		}
	}
	WriteResponse(Http.StatusOK, response, request, productModels)
}

func (productWebHandler *productWebHandler) HandleGetProductV1(response Http.ResponseWriter, request *Http.Request) {
	var product Domain.Product
	var ID Types.ID
	var err error
	if ID, err = GetID(request); err != nil {
		BadRequest(response, request, BadRequestTitle, make(map[string]any))
		return
	}

	if product, err = productWebHandler.productService.FindByID(ID.ID); err != nil {
		StatusNotFound(response, request)
		return
	}
	WriteResponse(Http.StatusOK, response, request, product)
}

func ProductWebHandler(productService Service.IProductService, productMapper Mapper.ProductMapper, categoryMapper Mapper.ICategoryMapper) IProductWebHandler {
	productWebHandler := &productWebHandler{productService, productMapper, categoryMapper}
	return productWebHandler
}
