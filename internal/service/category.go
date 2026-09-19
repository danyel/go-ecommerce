package service

import (
	Domain "github.com/danyel/ecommerce/internal/domain"
	Model "github.com/danyel/ecommerce/internal/model"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Types "github.com/danyel/ecommerce/internal/types"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ICategoryService interface {
	GetCategories() []Domain.Category
	GetCategory(ID Uuid.UUID) (Domain.Category, error)
	CreateCategory(createCategory Model.CreateCategory) (Uuid.UUID, error)
}

type categoryService struct {
	categoryRepository Persistence.CrudRepository[Persistence.CategoryModel]
}

func (categoryService *categoryService) GetCategories() []Domain.Category {
	categories := categoryService.categoryRepository.FindAll(Persistence.SearchCriteria{Preloads: []string{"Children"}})
	return mapCategories(categories)
}

func (categoryService *categoryService) GetCategory(ID Uuid.UUID) (Domain.Category, error) {
	var category Domain.Category
	categoryModel, err := categoryService.categoryRepository.FindByID(ID)
	if err != nil {
		return category, err
	}
	return mapCategory(*categoryModel), err
}

func (categoryService *categoryService) CreateCategory(createCategory Model.CreateCategory) (Uuid.UUID, error) {
	var err error
	category := &Persistence.CategoryModel{
		Name: createCategory.Name,
	}

	if err := categoryService.categoryRepository.Create(category); err != nil {
		return Uuid.Nil, err
	}
	var children []*Persistence.CategoryModel
	if len(createCategory.Children) > 0 {
		children = categoryService.categoryRepository.FindAll(Persistence.SearchCriteria{
			WhereClause: Persistence.WhereClause{
				Query:  "id IN ?",
				Params: []any{createCategory.Children},
			},
		})
	}

	if len(children) > 0 {
		if err = categoryService.categoryRepository.AssocAppend(category, "Children", children); err != nil {
			return category.ID, err
		}
	}
	return category.ID, err
}

func mapCategories(categoryModels []*Persistence.CategoryModel) []Domain.Category {
	categories := make([]Domain.Category, len(categoryModels))

	for index, categoryModel := range categoryModels {
		categories[index] = Domain.Category{
			ID:   Types.NewID(categoryModel.ID),
			Name: categoryModel.Name,
			// Important: children as pointers
			Children: mapCategories(categoryModel.Children),
		}
	}

	return categories
}

func mapCategory(categoryModel Persistence.CategoryModel) Domain.Category {
	return Domain.Category{
		ID:   Types.NewID(categoryModel.ID),
		Name: categoryModel.Name,
		// Important: children as pointers
		Children: mapCategories(categoryModel.Children),
	}
}

func CategoryService(categoryRepository Persistence.CrudRepository[Persistence.CategoryModel]) ICategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
