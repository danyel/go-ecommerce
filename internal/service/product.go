// Package service the service layer
package service

import (
	Domain "github.com/danyel/ecommerce/internal/domain"
	Mapper "github.com/danyel/ecommerce/internal/mapper"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Port "github.com/danyel/ecommerce/internal/port"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type IProductService interface {
	FindAll() []Domain.Product
	FindPage(page int, pageSize int) ([]Domain.Product, int64)
	FindByID(uuid Uuid.UUID) (Domain.Product, error)
	Update(product Domain.Product) error
	UpdateStock(ID Uuid.UUID, shoppingBasketID Uuid.UUID, stock int) error
}

func (productService *productService) FindPage(page int, pageSize int) ([]Domain.Product, int64) {
	orderBy := "created_at asc"
	offset := (page - 1) * pageSize
	models, total := productService.productRepository.Paginate(Persistence.SearchCriteria{
		Limit: &pageSize, Offset: &offset, OrderBy: &orderBy,
	})
	return productService.productMapper.MapProducts(modelsToPointers(models)), total
}

func modelsToPointers(models []Persistence.ProductModel) []*Persistence.ProductModel {
	result := make([]*Persistence.ProductModel, len(models))
	for i := range models {
		result[i] = &models[i]
	}
	return result
}

type productService struct {
	productRepository Persistence.CrudRepository[Persistence.ProductModel]
	productMapper     Mapper.ProductMapper
	publisher         Port.EventPublisher
}

func (productService *productService) FindAll() []Domain.Product {
	orderBy := "created_at asc"
	products := productService.productRepository.FindAll(Persistence.SearchCriteria{OrderBy: &orderBy})
	return productService.productMapper.MapProducts(products)
}

func (productService *productService) FindByID(ID Uuid.UUID) (Domain.Product, error) {
	var product Domain.Product
	productModel, err := productService.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}
	return productService.productMapper.MapProduct(productModel), nil
}

func (productService *productService) UpdateStock(ID Uuid.UUID, shoppingBasketID Uuid.UUID, stock int) error {
	product, err := productService.productRepository.FindByID(ID)
	if err != nil {
		return err
	}
	if stock > 0 {
		product.Stock -= stock
	} else if stock < 0 {
		product.Stock += -stock
	}
	if shoppingBasketID != Uuid.Nil {
		// we have a shopping basket
		// A shopping basket ID is retained for event consumers.
	}

	return productService.productRepository.Update(product)
}

func (productService *productService) Update(product Domain.Product) error {
	productModel, err := productService.productRepository.FindByID(product.ID.ID)
	if err != nil {
		return err
	}
	updateFields(product, productModel)
	return productService.productRepository.Update(productModel)
}

func updateFields(source Domain.Product, target *Persistence.ProductModel) {
	target.Brand = source.Brand
	target.Name = source.Name
	target.Description = source.Description
	target.Code = source.Code
	target.Price = float64(source.Price.Inclusive)
	target.CategoryID = source.Category.ID.ID
	target.ImageURL = source.ImageURL
	target.Stock = source.Stock
}

func ProductService(productRepository Persistence.CrudRepository[Persistence.ProductModel], productMapper Mapper.ProductMapper, publisher Port.EventPublisher) IProductService {
	productService := &productService{productRepository: productRepository, productMapper: productMapper, publisher: publisher}
	return productService
}
