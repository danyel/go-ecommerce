package category

import (
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryService interface {
	GetCategories() []Category
	GetCategory(categoryID Uuid.UUID) (Category, error)
	CreateCategory(createCategory CreateCategory) (CategoryID, error)
}

type categoryService struct {
	categoryRepository CommonRepository.CrudRepository[CategoryModel]
}

func (s *categoryService) GetCategories() []Category {
	categoryModels := s.categoryRepository.FindAll(CommonRepository.SearchCriteria{Preloads: []string{"Children"}})
	return mapCategories(categoryModels)
}

func (s *categoryService) GetCategory(categoryID Uuid.UUID) (Category, error) {
	var category Category
	categoryModel, err := s.categoryRepository.FindById(categoryID)
	if err != nil {
		return category, err
	}
	return mapCategory(*categoryModel), err
}

func (s *categoryService) CreateCategory(createCategory CreateCategory) (CategoryID, error) {
	var err error
	var categoryID CategoryID
	category := &CategoryModel{
		Name: createCategory.Name,
	}

	if err := s.categoryRepository.Create(category); err != nil {
		return categoryID, err
	}
	var children []*CategoryModel
	if len(createCategory.Children) > 0 {
		children = s.categoryRepository.FindAll(CommonRepository.SearchCriteria{
			WhereClause: CommonRepository.WhereClause{
				Query:  "id IN ?",
				Params: []any{createCategory.Children},
			},
		})
	}

	if len(children) > 0 {
		if err = s.categoryRepository.AssocAppend(category, "Children", children); err != nil {
			return categoryID, err
		}
	}
	categoryID.ID = category.ID
	return categoryID, err
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

func NewCategoryService(DB *Database.DB) CategoryService {
	return &categoryService{
		categoryRepository: CommonRepository.NewCrudRepository[CategoryModel](DB),
	}
}
