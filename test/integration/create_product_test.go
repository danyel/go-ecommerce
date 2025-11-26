package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dnoulet/ecommerce/internal/category"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/stretchr/testify/assert"
)

func Test_CreateProduct(t *testing.T) {
	app := SetupTestApp(t)

	categoryRepository := commonRepository.NewCrudRepository[category.CategoryModel](app.DB)
	t.Run("CreateProduct", func(t *testing.T) {
		c := category.CategoryModel{Name: "test", Children: []*category.CategoryModel{}}
		e := categoryRepository.Create(&c)
		assert.Nil(t, e)
		body := map[string]any{
			"brand":       "Apple",
			"name":        "iPhone 16",
			"description": "test device",
			"code":        "some code",
			"price":       10,
			"image_url":   "image_url",
			"categoryId":  c.ID.String(),
		}
		b, _ := json.Marshal(body)

		res, err := http.Post(app.Server.URL+"/api/product-management/v1/products", "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("HTTP request failed: %v", err)
		}

		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", res.StatusCode)
		}
	})
}
