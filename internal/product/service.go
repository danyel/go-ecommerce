package product

import (
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Port "github.com/danyel/ecommerce/internal/common/port"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService interface {
	FindAll() []Product
	FindById(uuid Uuid.UUID) (Product, error)
	Update(product Product) error
	UpdateStock(ID Uuid.UUID, shoppingBasketId Uuid.UUID, stock int) error
}

type productService struct {
	productRepository Repository.CrudRepository[ProductModel]
	productMapper     ProductMapper
	publisher         Port.EventPublisher
}

func (productService *productService) FindAll() []Product {
	orderBy := "created_at asc"
	products := productService.productRepository.FindAll(Repository.SearchCriteria{OrderBy: &orderBy})
	return productService.productMapper.MapProducts(products)
}

func (productService *productService) FindById(ID Uuid.UUID) (Product, error) {
	var product Product
	productModel, err := productService.productRepository.FindById(ID)
	if err != nil {
		return product, err
	}
	Logger.Log.Debug("Getting model: %v", productModel)

	return productService.productMapper.MapProduct(productModel), nil
}

func (productService *productService) UpdateStock(ID Uuid.UUID, shoppingBasketId Uuid.UUID, stock int) error {
	product, err := productService.productRepository.FindById(ID)
	if err != nil {
		return err
	}
	if stock > 0 {
		product.Stock -= stock
	} else if stock < 0 {
		product.Stock += -stock
	}
	if shoppingBasketId != Uuid.Nil {

	}

	return productService.productRepository.Update(product)
}

func (productService *productService) Update(product Product) error {
	productModel, err := productService.productRepository.FindById(product.ID.ID)
	if err != nil {
		return err
	}
	updateFields(product, productModel)
	return productService.productRepository.Update(productModel)
}

func updateFields(source Product, target *ProductModel) {
	target.Brand = source.Brand
	target.Name = source.Name
	target.Description = source.Description
	target.Code = source.Code
	target.Price = float64(source.Price.Inclusive)
	target.CategoryID = source.Category.ID
	target.ImageURL = source.ImageURL
	target.Stock = source.Stock
}

func NewService(productRepository Repository.CrudRepository[ProductModel], productMapper ProductMapper, publisher Port.EventPublisher) ProductService {
	productService := &productService{productRepository, productMapper, publisher}
	return productService
}
