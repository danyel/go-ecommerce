package product

import (
	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
)

type ProductMapper interface {
	MapProducts(models []*ProductModel) []Product
	MapProduct(productModel *ProductModel) Product
}

type productMapper struct {
	categoryService category.CategoryService
	mapCategory     func(model category.Category) category.Category
	cmsService      cms.CmsService
}

func (p *productMapper) MapProducts(models []*ProductModel) []Product {
	result := make([]Product, len(models))
	for i, productModel := range models {
		result[i] = p.MapProduct(productModel)
	}
	return result
}
func (p *productMapper) MapProduct(productModel *ProductModel) Product {
	categoryModel, _ := p.categoryService.GetCategory(productModel.CategoryId)
	description, _ := p.cmsService.GetTranslation(productModel.Description, "nl_BE")
	name, _ := p.cmsService.GetTranslation(productModel.Name, "nl_BE")
	return Product{
		Code:        productModel.Code,
		Price:       productModel.Price,
		Category:    categoryModel,
		ImageUrl:    productModel.ImageUrl,
		Brand:       productModel.Brand,
		Description: description.Value,
		Name:        name.Value,
		ID:          productModel.ID,
		Stock:       productModel.Stock,
	}
}
func NewProductMapper(p category.CategoryService, c cms.CmsService) ProductMapper {
	return &productMapper{categoryService: p, cmsService: c}
}
