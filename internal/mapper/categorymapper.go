// Package mapper: mapper layer
package mapper

import (
	Domain "github.com/danyel/ecommerce/internal/domain"
	Model "github.com/danyel/ecommerce/internal/model"
)

type ICategoryMapper interface {
	Map(category Domain.Category) Model.CategoryDTO
}

type categoryMapper struct {
}

func (categoryMapper *categoryMapper) Map(category Domain.Category) Model.CategoryDTO {
	return Model.CategoryDTO{
		ID:       category.ID,
		Name:     category.Name,
		Children: mapCategories(category.Children),
	}
}

func mapCategories(categories []Domain.Category) []Model.CategoryDTO {
	result := make([]Model.CategoryDTO, len(categories))
	for i, category := range categories {
		result[i] = (&categoryMapper{}).Map(category)
	}
	return result
}

func CategoryMapper() ICategoryMapper {
	categoryMapper := &categoryMapper{}
	return categoryMapper
}
