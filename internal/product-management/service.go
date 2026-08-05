package product_management

import (
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"github.com/danyel/ecommerce/internal/common/types"
	"github.com/danyel/ecommerce/internal/product"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService interface {
	GetProducts() []product.Product
	GetProduct(id types.Id) (product.Product, error)
	DeleteProduct(id types.Id) error
	UpdateProduct(id types.Id, updateProduct UpdateProduct) error
	CreateProduct(createProduct CreateProduct) (ProductId, error)
}

type productService struct {
	productRepository commonRepository.CrudRepository[product.ProductModel]
	productService    product.ProductService
}

func (s *productService) GetProducts() []product.Product {
	return s.productService.GetProducts()
}

func (s *productService) GetProduct(uuid types.Id) (product.Product, error) {
	var p product.Product
	productModel, err := s.productService.GetProduct(uuid.ID)
	if err != nil {
		return p, err
	}
	return product.Product{
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
func (s *productService) DeleteProduct(uuid types.Id) error {
	return s.productRepository.Delete(uuid.ID)
}
func (s *productService) UpdateProduct(uuid types.Id, updateProduct UpdateProduct) error {
	p, err := s.productRepository.FindById(uuid.ID)
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
func NewProductService(DB *gorm.DB) ProductService {
	return &productService{
		commonRepository.NewCrudRepository[product.ProductModel](DB),
		product.NewProductService(DB),
	}
}
