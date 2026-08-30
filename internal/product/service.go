package product

import (
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService interface {
	GetProducts() []Product
	GetProduct(uuid Uuid.UUID) (Product, error)
}

type productService struct {
	productRepository Repository.CrudRepository[ProductModel]
	productMapper     ProductMapper
}

func (productService *productService) GetProducts() []Product {
	orderBy := "created_at asc"
	products := productService.productRepository.FindAll(Repository.SearchCriteria{OrderBy: &orderBy})
	return productService.productMapper.MapProducts(products)
}

func (productService *productService) GetProduct(ID Uuid.UUID) (Product, error) {
	var product Product
	productModel, err := productService.productRepository.FindByID(ID)
	if err != nil {
		return product, err
	}

	return productService.productMapper.MapProduct(productModel), nil
}

func NewService(productRepository Repository.CrudRepository[ProductModel], productMapper ProductMapper) ProductService {
	productService := &productService{productRepository, productMapper}
	return productService
}
