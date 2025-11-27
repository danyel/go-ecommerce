package product_management

import (
	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	"github.com/danyel/ecommerce/internal/product"
)

type CategoryMapper interface {
	MapCategories(models []category.Category) []Category
	MapCategory(model category.Category) Category
}

type ProductMapper interface {
	MapProducts(models []product.Product) []Product
	MapProduct(productModel product.Product) Product
}

type productMapper struct {
	categoryService category.CategoryService
	mapCategory     func(model category.Category) Category
	cmsService      cms.CmsService
}

type categoryMapper struct{}

func (p *productMapper) MapProducts(models []product.Product) []Product {
	result := make([]Product, len(models))
	for i, productModel := range models {
		result[i] = p.MapProduct(productModel)
	}
	return result
}
func (p *productMapper) MapProduct(productModel product.Product) Product {
	categoryModel, _ := p.categoryService.GetCategory(productModel.CategoryId)
	description, _ := p.cmsService.GetTranslation(productModel.Description, "nl_BE")
	name, _ := p.cmsService.GetTranslation(productModel.Name, "nl_BE")
	return Product{
		Code:        productModel.Code,
		Price:       productModel.Price,
		Category:    p.mapCategory(categoryModel),
		ImageUrl:    productModel.ImageUrl,
		Brand:       productModel.Brand,
		Description: description.Value,
		Name:        name.Value,
		ID:          productModel.ID,
		Stock:       productModel.Stock,
	}
}

func (c *categoryMapper) MapCategories(models []category.Category) []Category {
	categories := make([]Category, len(models))

	for i, m := range models {
		categories[i] = Category{
			ID:       m.ID,
			Name:     m.Name,
			Children: c.MapCategories(m.Children),
		}
	}

	return categories
}

func (c *categoryMapper) MapCategory(model category.Category) Category {
	return Category{
		ID:       model.ID,
		Name:     model.Name,
		Children: c.MapCategories(model.Children),
	}
}

func NewProductMapper(p category.CategoryService, c cms.CmsService) ProductMapper {
	return &productMapper{categoryService: p, mapCategory: NewCategoryMapper().MapCategory, cmsService: c}
}

func NewCategoryMapper() CategoryMapper {
	return &categoryMapper{}
}
