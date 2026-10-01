package handler

import (
	Fmt "fmt"
	Http "net/http"
	Strconv "strconv"

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
	maxPageSize    int
}

func (productWebHandler *productWebHandler) HandleGetProductsV1(response Http.ResponseWriter, request *Http.Request) {
	page, pageSize, err := ParsePagination(request, productWebHandler.maxPageSize)
	if err != nil {
		BadRequest(response, request, InvalidRequestTitle, map[string]any{"pagination": err.Error()})
		return
	}
	products, total := productWebHandler.productService.FindPage(page, pageSize)
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
	WriteResponse(Http.StatusOK, response, request, Model.Page[Model.ProductDTO]{
		Items: productModels, Page: page, PageSize: pageSize, Total: int(total),
		HasNextPage: page*pageSize < int(total),
	})
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

func ProductWebHandler(productService Service.IProductService, productMapper Mapper.ProductMapper, categoryMapper Mapper.ICategoryMapper, maxPageSize ...int) IProductWebHandler {
	maximum := 50
	if len(maxPageSize) > 0 && maxPageSize[0] > 0 {
		maximum = maxPageSize[0]
	}
	productWebHandler := &productWebHandler{productService, productMapper, categoryMapper, maximum}
	return productWebHandler
}

func ParsePagination(request *Http.Request, maxPageSize int) (int, int, error) {
	page, pageSize := 1, maxPageSize
	var err error
	if value := request.URL.Query().Get("page"); value != "" {
		page, err = Strconv.Atoi(value)
		if err != nil || page < 1 {
			return 0, 0, Fmt.Errorf("page must be a positive integer")
		}
	}
	if value := request.URL.Query().Get("page_size"); value != "" {
		pageSize, err = Strconv.Atoi(value)
		if err != nil || pageSize < 1 || pageSize > maxPageSize {
			return 0, 0, Fmt.Errorf("page_size must be between 1 and %d", maxPageSize)
		}
	}
	return page, pageSize, nil
}
