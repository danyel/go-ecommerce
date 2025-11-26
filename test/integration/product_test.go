package integration_test

import (
	"testing"

	"github.com/dnoulet/ecommerce/internal/category"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/internal/product"
	"github.com/dnoulet/ecommerce/test/integration"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	"github.com/stretchr/testify/assert"
)

func TestProductBackend(t *testing.T) {
	initializer := integration.NewBackendInitializer()
	initializer.Run()
	productRepository := commonRepository.NewCrudRepository[product.ProductModel](initializer.DB)
	categoryRepository := commonRepository.NewCrudRepository[category.CategoryModel](initializer.DB)

	t.Run("Create Product ", func(t *testing.T) {
		c := category.CategoryModel{Name: "test", Children: []*category.CategoryModel{}}
		e := categoryRepository.Create(&c)
		assert.Nil(t, e)
		p := product.ProductModel{
			Brand:       "Brand",
			Name:        "Name",
			Description: "Description",
			Code:        "Code",
			Price:       1000,
			CategoryId:  c.ID,
			ImageUrl:    "ImageUrl",
			Stock:       1,
		}

		err := productRepository.Create(&p)
		assert.Nil(t, err)
	})
}
