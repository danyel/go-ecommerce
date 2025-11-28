package product

import (
	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
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
	productMapper     ProductMapper
}

func (s *productService) GetProducts() []Product {
	orderBy := "created_at asc"
	products := s.productRepository.FindAll(commonRepository.SearchCriteria{OrderBy: &orderBy})
	return s.productMapper.MapProducts(products)
}
func (s *productService) GetProduct(uuid uuid.UUID) (Product, error) {
	var product Product
	productModel, err := s.productRepository.FindById(uuid)
	if err != nil {
		return product, err
	}

	return s.productMapper.MapProduct(productModel), nil
}

func NewProductService(DB *gorm.DB) ProductService {
	s := &productService{commonRepository.NewCrudRepository[ProductModel](DB), NewProductMapper(category.NewCategoryService(DB), cms.NewCmsService(DB))}
	return s
}
