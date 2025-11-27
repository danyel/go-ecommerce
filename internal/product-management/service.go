package product_management

import (
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"github.com/danyel/ecommerce/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService interface {
	GetProducts() []product.Product
	GetProduct(uuid uuid.UUID) (product.Product, error)
	DeleteProduct(uuid uuid.UUID) error
	UpdateProduct(uuid uuid.UUID, updateProduct UpdateProduct) error
	CreateProduct(createProduct CreateProduct) (ProductId, error)
}

type productService struct {
	productRepository commonRepository.CrudRepository[product.ProductModel]
}

func (s *productService) GetProducts() []product.Product {
	orderBy := "created_at asc"
	products := s.productRepository.FindAll(commonRepository.SearchCriteria{OrderBy: &orderBy})
	return MapToProduct(products)
}

func (s *productService) GetProduct(uuid uuid.UUID) (product.Product, error) {
	var p product.Product
	productModel, err := s.productRepository.FindById(uuid)
	if err != nil {
		return p, err
	}
	return product.Product{
		Code:        productModel.Code,
		Price:       productModel.Price,
		Stock:       productModel.Stock,
		CategoryId:  productModel.CategoryId,
		ImageUrl:    productModel.ImageUrl,
		Brand:       productModel.Brand,
		Description: productModel.Description,
		Name:        productModel.Name,
		ID:          productModel.ID,
	}, nil
}
func (s *productService) DeleteProduct(uuid uuid.UUID) error {
	return s.productRepository.Delete(uuid)
}
func (s *productService) UpdateProduct(uuid uuid.UUID, updateProduct UpdateProduct) error {
	p, err := s.productRepository.FindById(uuid)
	if err != nil {
		return err
	}
	p.Name = updateProduct.Name
	p.Brand = updateProduct.Brand
	p.Description = updateProduct.Description
	p.Stock = updateProduct.Stock
	p.CategoryId = updateProduct.CategoryId
	p.ImageUrl = updateProduct.ImageUrl
	p.Price = updateProduct.Price
	return s.productRepository.Update(p)
}
func (s *productService) CreateProduct(createProduct CreateProduct) (ProductId, error) {
	var productId ProductId
	p := product.ProductModel{
		Code:        createProduct.Code,
		Price:       createProduct.Price,
		CategoryId:  createProduct.CategoryId,
		ImageUrl:    createProduct.ImageUrl,
		Brand:       createProduct.Brand,
		Description: createProduct.Description,
		Name:        createProduct.Name,
	}
	err := s.productRepository.Create(&p)
	if err != nil {
		return productId, err
	}
	return ProductId{ID: p.ID}, nil
}

func MapToProduct(productModels []*product.ProductModel) []product.Product {
	result := make([]product.Product, len(productModels))
	for i, p := range productModels {
		result[i] = product.Product{
			Code:        p.Code,
			Price:       p.Price,
			CategoryId:  p.CategoryId,
			ImageUrl:    p.ImageUrl,
			Brand:       p.Brand,
			Description: p.Description,
			Name:        p.Name,
			ID:          p.ID,
			Stock:       p.Stock,
		}
	}
	return result
}

func NewProductService(DB *gorm.DB) ProductService {
	return &productService{
		commonRepository.NewCrudRepository[product.ProductModel](DB),
	}
}
