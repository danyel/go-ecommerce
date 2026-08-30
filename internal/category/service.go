package category

import (
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryService interface {
	GetCategories() []Category
	GetCategory(ID Uuid.UUID) (Category, error)
	CreateCategory(createCategory CreateCategory) (Uuid.UUID, error)
}

type categoryService struct {
	categoryRepository Repository.CrudRepository[CategoryModel]
}

func (categoryService *categoryService) GetCategories() []Category {
	categories := categoryService.categoryRepository.FindAll(Repository.SearchCriteria{Preloads: []string{"Children"}})
	return mapCategories(categories)
}

func (categoryService *categoryService) GetCategory(ID Uuid.UUID) (Category, error) {
	var category Category
	categoryModel, err := categoryService.categoryRepository.FindByID(ID)
	if err != nil {
		return category, err
	}
	return mapCategory(*categoryModel), err
}

func (categoryService *categoryService) CreateCategory(createCategory CreateCategory) (Uuid.UUID, error) {
	var err error
	category := &CategoryModel{
		Name: createCategory.Name,
	}

	if err := categoryService.categoryRepository.Create(category); err != nil {
		return Uuid.Nil, err
	}
	var children []*CategoryModel
	if len(createCategory.Children) > 0 {
		children = categoryService.categoryRepository.FindAll(Repository.SearchCriteria{
			WhereClause: Repository.WhereClause{
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

func mapCategories(categoryModels []*CategoryModel) []Category {
	categories := make([]Category, len(categoryModels))

	for index, categoryModel := range categoryModels {
		categories[index] = Category{
			ID:   Types.NewID(categoryModel.ID),
			Name: categoryModel.Name,
			// Important: children as pointers
			Children: mapCategories(categoryModel.Children),
		}
	}

	return categories
}

func mapCategory(categoryModel CategoryModel) Category {
	return Category{
		ID:   Types.NewID(categoryModel.ID),
		Name: categoryModel.Name,
		// Important: children as pointers
		Children: mapCategories(categoryModel.Children),
	}
}

func NewService(categoryRepository Repository.CrudRepository[CategoryModel]) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
