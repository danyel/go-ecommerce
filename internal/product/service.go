package product

import (
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService interface {
	GetProducts() []Product
	GetProduct(uuid Uuid.UUID) (Product, error)
}

type productService struct {
	productRepository CommonRepository.CrudRepository[ProductModel]
	productMapper     ProductMapper
}

func (s *productService) GetProducts() []Product {
	orderBy := "created_at asc"
	products := s.productRepository.FindAll(CommonRepository.SearchCriteria{OrderBy: &orderBy})
	return s.productMapper.MapProducts(products)
}

func (s *productService) GetProduct(uuid Uuid.UUID) (Product, error) {
	var product Product
	productModel, err := s.productRepository.FindById(uuid)
	if err != nil {
		return product, err
	}

	return s.productMapper.MapProduct(productModel), nil
}

func NewProductService(DB *Database.DB) ProductService {
	s := &productService{CommonRepository.NewCrudRepository[ProductModel](DB), NewProductMapper(Category.NewCategoryService(DB), CMS.NewCmsService(DB))}
	return s
}
