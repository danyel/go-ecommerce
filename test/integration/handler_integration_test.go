package integration

import (
	"testing"

	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/internal/management"
	productmanagement "github.com/dnoulet/ecommerce/internal/product-management"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	app := SetupWebIntegration(t)

	t.Run("Product Handler", func(t *testing.T) {
		categoryRepository := commonRepository.NewCrudRepository[category.CategoryModel](app.DB)
		t.Run("CreateProduct", func(t *testing.T) {
			c := category.CategoryModel{Name: "test", Children: []*category.CategoryModel{}}
			e := categoryRepository.Create(&c)
			assert.Nil(t, e)
			body := &productmanagement.CreateProduct{
				Brand:       "Apple",
				Name:        "iPhone 16",
				Description: "test device",
				Code:        "some code",
				Price:       10,
				ImageUrl:    "image_url",
				CategoryId:  c.ID,
			}
			var productId productmanagement.ProductId
			app.Post(app.Server.URL+"/api/product-management/v1/products", "application/json", body).
				GetResponseBody(&productId).
				AssertStatusCreated()
			assert.NotNil(t, productId.ID)
		})
	})

	t.Run("CMS Handler", func(t *testing.T) {
		f := NewFixture(commonRepository.NewCrudRepository[cms.CmsModel](app.DB))
		f.SaveModel(&cms.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"})
		f.SaveModel(&cms.CmsModel{Code: "code", Language: "nl_fr", Value: "Value_fr"})
		f.SaveModel(&cms.CmsModel{Code: "another_code", Language: "nl_be", Value: "AnotherValue_nl"})
		f.SaveModel(&cms.CmsModel{Code: "another_code", Language: "nl_fr", Value: "AnotherValue_fr"})
		f.SaveModel(&cms.CmsModel{Code: "yet_another_code", Language: "nl_be", Value: "YetAnotherValue_nl"})
		f.SaveModel(&cms.CmsModel{Code: "yet_another_code", Language: "nl_fr", Value: "YetAnotherValue_fr"})

		t.Run("TestCmsHandler", func(t *testing.T) {
			t.Run("CmsHandler retrieve dutch", func(t *testing.T) {
				var translations []cms.Translation
				app.Get(app.Server.URL + "/api/cms/v1/translations?language=nl_be").
					GetResponseBody(&translations).
					AssertStatusOk()
				assert.Equal(t, 3, len(translations))
				assert.Equal(t, "code", translations[0].Code)
				assert.Equal(t, "Value_nl", translations[0].Value)
				assert.Equal(t, "another_code", translations[1].Code)
				assert.Equal(t, "AnotherValue_nl", translations[1].Value)
				assert.Equal(t, "yet_another_code", translations[2].Code)
				assert.Equal(t, "YetAnotherValue_nl", translations[2].Value)
			})

			t.Run("CmsHandler retrieve french", func(t *testing.T) {
				var translations []cms.Translation
				app.Get(app.Server.URL + "/api/cms/v1/translations?language=nl_fr").
					GetResponseBody(&translations).
					AssertStatusOk()
				assert.Equal(t, 3, len(translations))
				assert.Equal(t, "code", translations[0].Code)
				assert.Equal(t, "Value_fr", translations[0].Value)
				assert.Equal(t, "another_code", translations[1].Code)
				assert.Equal(t, "AnotherValue_fr", translations[1].Value)
				assert.Equal(t, "yet_another_code", translations[2].Code)
				assert.Equal(t, "YetAnotherValue_fr", translations[2].Value)
			})

			t.Run("CmsHandler retrieve all", func(t *testing.T) {
				var translations []cms.Translation
				app.Get(app.Server.URL + "/api/cms/v1/translations").
					GetResponseBody(&translations).
					AssertStatusOk()
				assert.Equal(t, 6, len(translations))
				assert.Equal(t, "code", translations[0].Code)
				assert.Equal(t, "Value_nl", translations[0].Value)
				assert.Equal(t, "code", translations[1].Code)
				assert.Equal(t, "Value_fr", translations[1].Value)
				assert.Equal(t, "another_code", translations[2].Code)
				assert.Equal(t, "AnotherValue_nl", translations[2].Value)
				assert.Equal(t, "another_code", translations[3].Code)
				assert.Equal(t, "AnotherValue_fr", translations[3].Value)
				assert.Equal(t, "yet_another_code", translations[4].Code)
				assert.Equal(t, "YetAnotherValue_nl", translations[4].Value)
				assert.Equal(t, "yet_another_code", translations[5].Code)
				assert.Equal(t, "YetAnotherValue_fr", translations[5].Value)
			})

			t.Run("CmsHandler retrieve none because of invalid language", func(t *testing.T) {
				var translations []cms.Translation
				app.Get(app.Server.URL + "/api/cms/v1/translations?language=nl_de").
					GetResponseBody(&translations).
					AssertStatusOk()
				assert.Equal(t, 0, len(translations))
			})
		})
	})

	t.Run("Management Handler", func(t *testing.T) {
		t.Run("ManagementHandler: Create a new translation but it already exist so return 400", func(t *testing.T) {
			body := &management.CreateCms{
				Code:     "code",
				Language: "nl_be",
				Value:    "Value_nl",
			}
			app.Post(app.Server.URL+"/api/management/v1/translations", "application/json", body).
				AssertBadRequest()
		})

		t.Run("ManagementHandler: Create a new translation and 201 is return", func(t *testing.T) {
			body := &management.CreateCms{
				Code:     "unknown",
				Language: "nl_fr",
				Value:    "Value_fr",
			}
			app.Post(app.Server.URL+"/api/management/v1/translations", "application/json", body).
				AssertStatusCreated()
		})
	})
}
