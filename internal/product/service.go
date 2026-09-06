package product

import (
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

func updateFields(product Product, productModel *ProductModel) {
	productModel.Brand = product.Brand
	productModel.Name = product.Name
	productModel.Description = product.Description
	productModel.Code = product.Code
	productModel.Price = float64(product.Price.Inclusive)
	productModel.CategoryID = product.Category.ID
	productModel.ImageURL = product.ImageURL
	productModel.Stock = product.Stock
}

func NewService(productRepository Repository.CrudRepository[ProductModel], productMapper ProductMapper, publisher Port.EventPublisher) ProductService {
	productService := &productService{productRepository, productMapper, publisher}
	return productService
}
