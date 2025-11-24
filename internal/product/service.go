package product

import (
	repository "github.com/dnoulet/ecommerce/internal/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService struct {
	GetProducts func() []Product
	GetProduct  func(uuid uuid.UUID) (Product, error)
}

func MapToProduct(productModels []*ProductModel) []Product {
	result := make([]Product, len(productModels))
	for i, product := range productModels {
		result[i] = Product{
			Code:        product.Code,
			Price:       product.Price,
			Category:    product.Category,
			ImageUrl:    product.ImageUrl,
			Brand:       product.Brand,
			Description: product.Description,
			Name:        product.Name,
			ID:          product.ID,
			Stock:       product.Stock,
		}
	}
	return result
}

func NewProductService(DB *gorm.DB) ProductService {
	productRepository := repository.NewCrudRepository[ProductModel](DB)
	return ProductService{
		GetProducts: func() []Product {
			orderBy := "created_at asc"
			products := productRepository.FindAll(repository.SearchCriteria{OrderBy: &orderBy})
			return MapToProduct(products)
		},
		GetProduct: func(uuid uuid.UUID) (Product, error) {
			var product Product
			productModel, err := productRepository.FindById(uuid)
			if err != nil {
				return product, err
			}
			return Product{
				Code:        productModel.Code,
				Price:       productModel.Price,
				Stock:       productModel.Stock,
				Category:    productModel.Category,
				ImageUrl:    productModel.ImageUrl,
				Brand:       productModel.Brand,
				Description: productModel.Description,
				Name:        productModel.Name,
				ID:          productModel.ID,
			}, nil
		},
	}
}
