package product

import (
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductMapper interface {
	MapProducts(models []*ProductModel) []Product
	MapProduct(productModel *ProductModel) Product
}

type productMapper struct {
	categoryService Category.CategoryService
	mapCategory     func(model Category.Category) Category.Category
	cmsService      CMS.CmsService
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
		Price:       Types.NewPrice(productModel.Price, "EUR"),
		Category:    categoryModel,
		ImageUrl:    productModel.ImageUrl,
		Brand:       productModel.Brand,
		Description: description.Value,
		Name:        name.Value,
		ID:          Types.NewID(productModel.ID),
		Stock:       productModel.Stock,
	}
}
func NewProductMapper(p Category.CategoryService, c CMS.CmsService) ProductMapper {
	return &productMapper{categoryService: p, cmsService: c}
}
