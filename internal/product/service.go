package product

import (
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService interface {
	GetProducts() []Product
	GetProduct(uuid uuid.UUID) (Product, error)
}

type productService struct {
	productRepository commonRepository.CrudRepository[ProductModel]
}

func (s *productService) GetProducts() []Product {
	orderBy := "created_at asc"
	products := s.productRepository.FindAll(commonRepository.SearchCriteria{OrderBy: &orderBy})
	return MapToProduct(products)
}
func (s *productService) GetProduct(uuid uuid.UUID) (Product, error) {
	var product Product
	productModel, err := s.productRepository.FindById(uuid)
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
	s := &productService{commonRepository.NewCrudRepository[ProductModel](DB)}
	return s
}
