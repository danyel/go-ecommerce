package product_management

import (
	repository "github.com/dnoulet/ecommerce/internal/common"
	"github.com/dnoulet/ecommerce/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductService struct {
	GetProducts   func() []product.Product
	GetProduct    func(uuid uuid.UUID) (product.Product, error)
	DeleteProduct func(uuid uuid.UUID) error
	UpdateProduct func(uuid uuid.UUID, updateProduct UpdateProduct) error
	CreateProduct func(createProduct CreateProduct) (ProductId, error)
}

func MapToProduct(productModels []*product.ProductModel) []product.Product {
	result := make([]product.Product, len(productModels))
	for i, p := range productModels {
		result[i] = product.Product{
			Code:        p.Code,
			Price:       p.Price,
			Category:    p.Category,
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
	productRepository := repository.NewCrudRepository[product.ProductModel](DB)
	return ProductService{
		GetProducts: func() []product.Product {
			orderBy := "created_at asc"
			products := productRepository.FindAll(repository.SearchCriteria{OrderBy: &orderBy})
			return MapToProduct(products)
		},
		GetProduct: func(uuid uuid.UUID) (product.Product, error) {
			var p product.Product
			productModel, err := productRepository.FindById(uuid)
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
		},
		DeleteProduct: func(uuid uuid.UUID) error {
			return productRepository.Delete(uuid)
		},
		UpdateProduct: func(uuid uuid.UUID, updateProduct UpdateProduct) error {
			p, err := productRepository.FindById(uuid)
			if err != nil {
				return err
			}
			p.Name = updateProduct.Name
			p.Brand = updateProduct.Brand
			p.Description = updateProduct.Description
			p.Stock = updateProduct.Stock
			p.Category = updateProduct.Category
			p.ImageUrl = updateProduct.ImageUrl
			p.Price = updateProduct.Price
			return productRepository.Update(p)
		},
		CreateProduct: func(createProduct CreateProduct) (ProductId, error) {
			var productId ProductId
			p := product.ProductModel{
				Code:        createProduct.Code,
				Price:       createProduct.Price,
				Category:    createProduct.Category,
				ImageUrl:    createProduct.ImageUrl,
				Brand:       createProduct.Brand,
				Description: createProduct.Description,
				Name:        createProduct.Name,
			}
			err := productRepository.Create(&p)
			if err != nil {
				return productId, err
			}
			return ProductId{ID: p.ID}, nil
		},
	}
}
