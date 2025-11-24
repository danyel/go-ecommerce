package product_management

import (
	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/product"
)

type CategoryMapper struct {
	MapCategories func(models []category.Category) []Category
	MapCategory   func(model category.Category) Category
}

type ProductMapper struct {
	MapProducts func(models []product.Product) []Product
	MapProduct  func(productModel product.Product) Product
}

func MapCategories(models []category.Category) []Category {
	categories := make([]Category, len(models))

	for i, m := range models {
		categories[i] = Category{
			ID:       m.ID,
			Name:     m.Name,
			Children: MapCategories(m.Children),
		}
	}

	return categories
}

func MapCategory(model category.Category) Category {
	return Category{
		ID:       model.ID,
		Name:     model.Name,
		Children: MapCategories(model.Children),
	}
}

func NewProductMapper(p *category.CategoryService) ProductMapper {
	return ProductMapper{
		MapProducts: func(models []product.Product) []Product {
			result := make([]Product, len(models))
			for i, productModel := range models {
				categoryModel, _ := p.GetCategory(productModel.Category)
				result[i] = Product{
					Code:        productModel.Code,
					Price:       productModel.Price,
					Category:    MapCategory(categoryModel),
					ImageUrl:    productModel.ImageUrl,
					Brand:       productModel.Brand,
					Description: productModel.Description,
					Name:        productModel.Name,
					ID:          productModel.ID,
					Stock:       productModel.Stock,
				}
			}
			return result
		},
		MapProduct: func(productModel product.Product) Product {
			categoryModel, _ := p.GetCategory(productModel.Category)
			return Product{
				Code:        productModel.Code,
				Price:       productModel.Price,
				Category:    MapCategory(categoryModel),
				ImageUrl:    productModel.ImageUrl,
				Brand:       productModel.Brand,
				Description: productModel.Description,
				Name:        productModel.Name,
				ID:          productModel.ID,
				Stock:       productModel.Stock,
			}
		},
	}
}

func NewCategoryMapper() CategoryMapper {
	return CategoryMapper{
		MapCategories: MapCategories,
		MapCategory:   MapCategory,
	}
}
