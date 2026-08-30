package productmanagement

import (
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductManagementService interface {
	GetProducts() []Product.Product
	GetProduct(id Types.ID) (Product.Product, error)
	DeleteProduct(id Types.ID) error
	UpdateProduct(id Types.ID, updateProduct Product.UpdateProductDTO) error
	CreateProduct(createProduct Product.CreateProductDTO) (Types.ID, error)
}

type productManagementService struct {
	productRepository Repository.CrudRepository[Product.ProductModel]
	productService    Product.ProductService
}

func (productManagementService *productManagementService) GetProducts() []Product.Product {
	return productManagementService.productService.GetProducts()
}

func (productManagementService *productManagementService) GetProduct(ID Types.ID) (Product.Product, error) {
	product, err := productManagementService.productService.GetProduct(ID.ID)
	if err != nil {
		return product, err
	}
	return Product.Product{
		Code:        product.Code,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		ImageURL:    product.ImageURL,
		Brand:       product.Brand,
		Description: product.Description,
		Name:        product.Name,
		ID:          product.ID,
	}, nil
}

func (productManagementService *productManagementService) DeleteProduct(ID Types.ID) error {
	return productManagementService.productRepository.Delete(ID.ID)
}

func (productManagementService *productManagementService) UpdateProduct(ID Types.ID, updateProduct Product.UpdateProductDTO) error {
	productModel, err := productManagementService.productRepository.FindByID(ID.ID)
	if err != nil {
		return err
	}
	productModel.Name = updateProduct.Name
	productModel.Brand = updateProduct.Brand
	productModel.Description = updateProduct.Description
	productModel.Stock = updateProduct.Stock
	productModel.CategoryID = updateProduct.CategoryID.ID
	productModel.ImageURL = updateProduct.ImageURL
	productModel.Price = updateProduct.Price
	return productManagementService.productRepository.Update(productModel)
}

func (productManagementService *productManagementService) CreateProduct(createProduct Product.CreateProductDTO) (Types.ID, error) {
	var productID Types.ID
	productModel := Product.ProductModel{
		Code:        createProduct.Code,
		Price:       createProduct.Price,
		CategoryID:  createProduct.CategoryID.ID,
		ImageURL:    createProduct.ImageURL,
		Brand:       createProduct.Brand,
		Description: createProduct.Description,
		Name:        createProduct.Name,
	}
	err := productManagementService.productRepository.Create(&productModel)
	if err != nil {
		return productID, err
	}
	return Types.NewID(productModel.ID), nil
}

func NewService(productRepository Repository.CrudRepository[Product.ProductModel], productService Product.ProductService) ProductManagementService {
	return &productManagementService{
		productRepository,
		productService,
	}
}
