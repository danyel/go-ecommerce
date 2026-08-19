package productmanagement

import (
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService interface {
	GetProducts() []Product.Product
	GetProduct(id Types.ID) (Product.Product, error)
	DeleteProduct(id Types.ID) error
	UpdateProduct(id Types.ID, updateProduct Product.UpdateProduct) error
	CreateProduct(createProduct Product.CreateProduct) (Types.ID, error)
}

type productService struct {
	productRepository CommonRepository.CrudRepository[Product.ProductModel]
	productService    Product.ProductService
}

func (s *productService) GetProducts() []Product.Product {
	return s.productService.GetProducts()
}

func (s *productService) GetProduct(uuid Types.ID) (Product.Product, error) {
	var p Product.Product
	productModel, err := s.productService.GetProduct(uuid.ID)
	if err != nil {
		return p, err
	}
	return Product.Product{
		Code:        productModel.Code,
		Price:       productModel.Price,
		Stock:       productModel.Stock,
		Category:    productModel.Category,
		ImageURL:    productModel.ImageURL,
		Brand:       productModel.Brand,
		Description: productModel.Description,
		Name:        productModel.Name,
		ID:          productModel.ID,
	}, nil
}

func (s *productService) DeleteProduct(uuid Types.ID) error {
	return s.productRepository.Delete(uuid.ID)
}

func (s *productService) UpdateProduct(uuid Types.ID, updateProduct Product.UpdateProduct) error {
	p, err := s.productRepository.FindById(uuid.ID)
	if err != nil {
		return err
	}
	p.Name = updateProduct.Name
	p.Brand = updateProduct.Brand
	p.Description = updateProduct.Description
	p.Stock = updateProduct.Stock
	p.CategoryID = updateProduct.CategoryID.ID
	p.ImageURL = updateProduct.ImageURL
	p.Price = updateProduct.Price
	return s.productRepository.Update(p)
}

func (s *productService) CreateProduct(createProduct Product.CreateProduct) (Types.ID, error) {
	var productID Types.ID
	p := Product.ProductModel{
		Code:        createProduct.Code,
		Price:       createProduct.Price,
		CategoryID:  createProduct.CategoryID.ID,
		ImageURL:    createProduct.ImageURL,
		Brand:       createProduct.Brand,
		Description: createProduct.Description,
		Name:        createProduct.Name,
	}
	err := s.productRepository.Create(&p)
	if err != nil {
		return productID, err
	}
	return Types.NewID(p.ID), nil
}

func NewProductService(DB *Database.DB) ProductService {
	return &productService{
		CommonRepository.NewCrudRepository[Product.ProductModel](DB),
		Product.NewProductService(DB),
	}
}
