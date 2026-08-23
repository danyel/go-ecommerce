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
	cmsService      CMS.CmsService
}

func (productMapper *productMapper) MapProducts(productModels []*ProductModel) []Product {
	result := make([]Product, len(productModels))
	for index, productModel := range productModels {
		result[index] = productMapper.MapProduct(productModel)
	}
	return result
}

func (productMapper *productMapper) MapProduct(productModel *ProductModel) Product {
	categoryModel, _ := productMapper.categoryService.GetCategory(productModel.CategoryID)
	// TODO fetch that information from the header
	description, _ := productMapper.cmsService.GetTranslation(productModel.Description, "nl_BE")
	// TODO fetch that information from the header
	name, _ := productMapper.cmsService.GetTranslation(productModel.Name, "nl_BE")
	return Product{
		Code:        productModel.Code,
		Price:       Types.NewPrice(productModel.Price, "EUR"),
		Category:    categoryModel,
		ImageURL:    productModel.ImageURL,
		Brand:       productModel.Brand,
		Description: description.Value,
		Name:        name.Value,
		ID:          Types.NewID(productModel.ID),
		Stock:       productModel.Stock,
	}
}

func NewProductMapper(categoryService Category.CategoryService, cmsService CMS.CmsService) ProductMapper {
	return &productMapper{categoryService: categoryService, cmsService: cmsService}
}
