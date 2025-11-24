package management

type CategoryMapper struct {
	MapCategories func(models []*CategoryModel) []Category
	MapCategory   func(model *CategoryModel) Category
}

func MapCategories(models []*CategoryModel) []Category {
	categories := make([]Category, len(models))

	for i, m := range models {
		categories[i] = Category{
			ID:   m.ID,
			Name: m.Name,
			// Important: children as pointers
			Children: MapCategories(m.Children),
		}
	}

	return categories
}

func MapCategory(model *CategoryModel) Category {
	return Category{
		ID:       model.ID,
		Name:     model.Name,
		Children: MapCategories(model.Children),
	}
}

func NewCategoryMapper() CategoryMapper {
	return CategoryMapper{
		MapCategories: MapCategories,
		MapCategory:   MapCategory,
	}
}
