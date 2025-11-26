package integration_test

import (
	"testing"

	"github.com/dnoulet/ecommerce/internal/category"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/test/integration"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	"github.com/stretchr/testify/assert"
)

func TestCategoryBackend(t *testing.T) {
	initializer := integration.NewBackendInitializer()
	initializer.Run()
	categoryRepository := commonRepository.NewCrudRepository[category.CategoryModel](initializer.DB)

	t.Run("Get Categories Before Creation", func(t *testing.T) {
		findAll := categoryRepository.FindAll(commonRepository.SearchCriteria{})
		assert.Equal(t, 0, len(findAll))
	})

	t.Run("Create Category", func(t *testing.T) {
		c := category.CategoryModel{Name: "test", Children: []*category.CategoryModel{}}
		e := categoryRepository.Create(&c)
		assert.Nil(t, e)
	})

	t.Run("Get Categories After Creation", func(t *testing.T) {
		findAll := categoryRepository.FindAll(commonRepository.SearchCriteria{})
		assert.Equal(t, 1, len(findAll))
	})
}
