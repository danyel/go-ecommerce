package service

import (
	Domain "github.com/danyel/ecommerce/internal/domain"
	Model "github.com/danyel/ecommerce/internal/model"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Types "github.com/danyel/ecommerce/internal/types"
)

//goland:noinspection GoNameStartsWithPackageName
type IProductManagementService interface {
	GetProducts() []Domain.Product
	GetProduct(id Types.ID) (Domain.Product, error)
	DeleteProduct(id Types.ID) error
	// UpdateProduct should be Domain.Product
	UpdateProduct(id Types.ID, updateProduct Model.UpdateProductDTO) error
	CreateProduct(createProduct Model.CreateProductDTO) (Types.ID, error)
}

type productManagementService struct {
	productRepository Persistence.CrudRepository[Persistence.ProductModel]
	productService    IProductService
}

func (productManagementService *productManagementService) GetProducts() []Domain.Product {
	return productManagementService.productService.FindAll()
}

func (productManagementService *productManagementService) GetProduct(ID Types.ID) (Domain.Product, error) {
	product, err := productManagementService.productService.FindByID(ID.ID)
	if err != nil {
		return product, err
	}
	return Domain.Product{
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

func (productManagementService *productManagementService) UpdateProduct(ID Types.ID, updateProduct Model.UpdateProductDTO) error {
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

func (productManagementService *productManagementService) CreateProduct(createProduct Model.CreateProductDTO) (Types.ID, error) {
	var productID Types.ID
	productModel := Persistence.ProductModel{
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

func ProductManagementService(productRepository Persistence.CrudRepository[Persistence.ProductModel], productService IProductService) IProductManagementService {
	return &productManagementService{
		productRepository,
		productService,
	}
}
