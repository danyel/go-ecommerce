package integration

import (
	"testing"

	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/internal/management"
	productmanagement "github.com/dnoulet/ecommerce/internal/product-management"
	"github.com/dnoulet/ecommerce/test/integration/initializer"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	wi := initializer.SetupWebIntegration(t)

	t.Run("Product Handler", func(t *testing.T) {
		categoryRepository := commonRepository.NewCrudRepository[category.CategoryModel](wi.Db())

		t.Run("CreateProduct", func(t *testing.T) {
			c := category.CategoryModel{Name: "test", Children: []*category.CategoryModel{}}
			assert.Nil(t, categoryRepository.Create(&c))
			b := &productmanagement.CreateProduct{
				Brand:       "Apple",
				Name:        "iPhone 16",
				Description: "test device",
				Code:        "some code",
				Price:       10,
				ImageUrl:    "image_url",
				CategoryId:  c.ID,
			}
			var productId productmanagement.ProductId
			wi.ProductManagementPostProducts(b).
				GetResponseBody(&productId).
				AssertStatusCreated().
				IsNotNil(productId.ID)
		})
	})

	t.Run("CMS Handler", func(t *testing.T) {
		f := NewFixture(commonRepository.NewCrudRepository[cms.CmsModel](wi.Db()))
		f.SaveModel(&cms.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"})
		f.SaveModel(&cms.CmsModel{Code: "code", Language: "nl_fr", Value: "Value_fr"})
		f.SaveModel(&cms.CmsModel{Code: "another_code", Language: "nl_be", Value: "AnotherValue_nl"})
		f.SaveModel(&cms.CmsModel{Code: "another_code", Language: "nl_fr", Value: "AnotherValue_fr"})
		f.SaveModel(&cms.CmsModel{Code: "yet_another_code", Language: "nl_be", Value: "YetAnotherValue_nl"})
		f.SaveModel(&cms.CmsModel{Code: "yet_another_code", Language: "nl_fr", Value: "YetAnotherValue_fr"})

		t.Run("TestCmsHandler", func(t *testing.T) {
			t.Run("CmsHandler retrieve dutch", func(t *testing.T) {
				var translations []cms.Translation
				wi.GetTranslations("nl_be").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(3, len(translations)).
					Equal("code", translations[0].Code).
					Equal("Value_nl", translations[0].Value).
					Equal("another_code", translations[1].Code).
					Equal("AnotherValue_nl", translations[1].Value).
					Equal("yet_another_code", translations[2].Code).
					Equal("YetAnotherValue_nl", translations[2].Value)
			})

			t.Run("CmsHandler retrieve french", func(t *testing.T) {
				var translations []cms.Translation
				wi.GetTranslations("nl_fr").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(3, len(translations)).
					Equal("code", translations[0].Code).
					Equal("Value_fr", translations[0].Value).
					Equal("another_code", translations[1].Code).
					Equal("AnotherValue_fr", translations[1].Value).
					Equal("yet_another_code", translations[2].Code).
					Equal("YetAnotherValue_fr", translations[2].Value)
			})

			t.Run("CmsHandler retrieve all", func(t *testing.T) {
				var translations []cms.Translation
				wi.GetTranslations("").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(6, len(translations)).
					Equal("code", translations[0].Code).
					Equal("Value_nl", translations[0].Value).
					Equal("code", translations[1].Code).
					Equal("Value_fr", translations[1].Value).
					Equal("another_code", translations[2].Code).
					Equal("AnotherValue_nl", translations[2].Value).
					Equal("another_code", translations[3].Code).
					Equal("AnotherValue_fr", translations[3].Value).
					Equal("yet_another_code", translations[4].Code).
					Equal("YetAnotherValue_nl", translations[4].Value).
					Equal("yet_another_code", translations[5].Code).
					Equal("YetAnotherValue_fr", translations[5].Value)
			})

			t.Run("CmsHandler retrieve none because of invalid language", func(t *testing.T) {
				var translations []cms.Translation
				wi.GetTranslations("nl_de").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(0, len(translations))
			})
		})
	})

	t.Run("Management Handler", func(t *testing.T) {
		t.Run("ManagementHandler: Create a new translation but it already exist so return 400", func(t *testing.T) {
			b := &management.CreateCms{
				Code:     "code",
				Language: "nl_be",
				Value:    "Value_nl",
			}
			wi.ManagementPostTranslations(b).
				AssertBadRequest()
		})

		t.Run("ManagementHandler: Create a new translation and 201 is return", func(t *testing.T) {
			b := &management.CreateCms{
				Code:     "unknown",
				Language: "nl_fr",
				Value:    "Value_fr",
			}
			var i management.CmsId
			wi.ManagementPostTranslations(b).
				GetResponseBody(&i).
				IsNotNil(i.ID).
				AssertStatusCreated()
		})
	})
}
