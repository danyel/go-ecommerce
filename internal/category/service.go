package category

import (
	repository "github.com/dnoulet/ecommerce/internal/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryService struct {
	GetCategories  func() []Category
	GetCategory    func(categoryID uuid.UUID) (Category, error)
	CreateCategory func(createCategory CreateCategory) (CategoryId, error)
}

func mapCategories(models []*CategoryModel) []Category {
	categories := make([]Category, len(models))

	for i, m := range models {
		categories[i] = Category{
			ID:   m.ID,
			Name: m.Name,
			// Important: children as pointers
			Children: mapCategories(m.Children),
		}
	}

	return categories
}

func mapCategory(categoryModel CategoryModel) Category {
	return Category{
		ID:   categoryModel.ID,
		Name: categoryModel.Name,
		// Important: children as pointers
		Children: mapCategories(categoryModel.Children),
	}
}

func NewCategoryService(DB *gorm.DB) CategoryService {
	categoryRepository := repository.NewCrudRepository[CategoryModel](DB)
	return CategoryService{
		GetCategories: func() []Category {
			categoryModels := categoryRepository.FindAll(repository.SearchCriteria{Preloads: []string{"Children"}})
			return mapCategories(categoryModels)
		},
		GetCategory: func(categoryID uuid.UUID) (Category, error) {
			var category Category
			categoryModel, err := categoryRepository.FindById(categoryID)
			if err != nil {
				return category, err
			}
			return mapCategory(*categoryModel), err
		},
		CreateCategory: func(createCategory CreateCategory) (CategoryId, error) {
			var err error
			var categoryId CategoryId
			category := &CategoryModel{
				Name: createCategory.Name,
			}

			if err := categoryRepository.Create(category); err != nil {
				return categoryId, err
			}
			var children []*CategoryModel
			if len(createCategory.Children) > 0 {
				children = categoryRepository.FindAll(repository.SearchCriteria{
					WhereClause: repository.WhereClause{
						Query:  "id IN ?",
						Params: []interface{}{createCategory.Children},
					},
				})
			}

			if len(children) > 0 {
				if err = categoryRepository.AssocAppend(category, createCategory.Name, createCategory.Children); err != nil {
					return categoryId, err
				}
			}
			categoryId.ID = category.ID
			return categoryId, err
		},
	}
}
