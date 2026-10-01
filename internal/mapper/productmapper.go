package mapper

import (
	Fmt "fmt"
	Strings "strings"

	Domain "github.com/danyel/ecommerce/internal/domain"
	Model "github.com/danyel/ecommerce/internal/model"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Types "github.com/danyel/ecommerce/internal/types"
	Uuid "github.com/google/uuid"
)

type CategoryService interface {
	GetCategory(ID Uuid.UUID) (Domain.Category, error)
}

type CmsService interface {
	GetTranslation(code string, language string) (Model.Translation, error)
}

//goland:noinspection GoNameStartsWithPackageName
type ProductMapper interface {
	MapProducts(models []*Persistence.ProductModel) []Domain.Product
	MapProduct(productModel *Persistence.ProductModel) Domain.Product
}

type productMapper struct {
	categoryService CategoryService
	cmsService      CmsService
}

func (productMapper *productMapper) MapProducts(productModels []*Persistence.ProductModel) []Domain.Product {
	result := make([]Domain.Product, len(productModels))
	for index, productModel := range productModels {
		result[index] = productMapper.MapProduct(productModel)
	}
	return result
}

func (productMapper *productMapper) MapProduct(productModel *Persistence.ProductModel) Domain.Product {
	prefix := Strings.Replace(productModel.Code, "-", "_", 1)
	categoryModel, _ := productMapper.categoryService.GetCategory(productModel.CategoryID)
	// TODO fetch that information from the header
	description, _ := productMapper.cmsService.GetTranslation(Fmt.Sprintf("%s_DESCRIPTION", prefix), "nl_BE")
	// TODO fetch that information from the header
	name, _ := productMapper.cmsService.GetTranslation(Fmt.Sprintf("%s_NAME", prefix), "nl_BE")
	return Domain.Product{
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

func NewProductMapper(categoryService CategoryService, cmsService CmsService) ProductMapper {
	return &productMapper{categoryService: categoryService, cmsService: cmsService}
}
