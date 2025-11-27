package integration

import (
	"testing"

	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/internal/product"
	"github.com/dnoulet/ecommerce/test/integration/initializer"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	"github.com/stretchr/testify/assert"
)

func TestServiceIntegration(t *testing.T) {
	bi := initializer.NewBackendInitializer()
	bi.Run()

	t.Run("Category Testing", func(t *testing.T) {
		categoryRepository := commonRepository.NewCrudRepository[category.CategoryModel](bi.Db())

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
	})

	t.Run("Cms Testing", func(t *testing.T) {
		cmsRepository := commonRepository.NewCrudRepository[cms.CmsModel](bi.Db())
		f := Database[cms.CmsModel](cmsRepository)
		cmsService := cms.NewCmsService(cmsRepository)

		t.Run("Create translation", func(t *testing.T) {
			e := cmsRepository.Create(&cms.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"})
			assert.Nil(t, e)
			e = cmsRepository.Create(&cms.CmsModel{Code: "code", Language: "nl_fr", Value: "Value_fr"})
			assert.Nil(t, e)
		})

		t.Run("Delete translation", func(t *testing.T) {
			p := &cms.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"}
			e := cmsRepository.Create(p)
			assert.Nil(t, e)
			err := cmsRepository.Delete(p.ID)
			assert.Nil(t, err)
		})

		t.Run("Get translations", func(t *testing.T) {
			f.Insert(&cms.CmsModel{Code: "another_code", Language: "nl_be", Value: "AnotherValue_nl"})
			f.Insert(&cms.CmsModel{Code: "another_code", Language: "nl_fr", Value: "AnotherValue_fr"})
			f.Insert(&cms.CmsModel{Code: "yet_another_code", Language: "nl_be", Value: "YetAnotherValue_nl"})
			f.Insert(&cms.CmsModel{Code: "yet_another_code", Language: "nl_fr", Value: "YetAnotherValue_fr"})

			t.Run("Given 3 codes in 2 languages When retrieving for dutch Then 3 codes have been returned", func(t *testing.T) {
				translations := cmsService.GetTranslations("nl_be")
				assert.Equal(t, 3, len(translations))
				assert.Equal(t, "code", translations[0].Code)
				assert.Equal(t, "Value_nl", translations[0].Value)
				assert.Equal(t, "another_code", translations[1].Code)
				assert.Equal(t, "AnotherValue_nl", translations[1].Value)
				assert.Equal(t, "yet_another_code", translations[2].Code)
				assert.Equal(t, "YetAnotherValue_nl", translations[2].Value)
			})

			t.Run("Given 3 codes in 2 languages When retrieving for french Then 3 codes have been returned", func(t *testing.T) {
				translations := cmsService.GetTranslations("nl_fr")
				assert.Equal(t, 3, len(translations))
				assert.Equal(t, "code", translations[0].Code)
				assert.Equal(t, "Value_fr", translations[0].Value)
				assert.Equal(t, "another_code", translations[1].Code)
				assert.Equal(t, "AnotherValue_fr", translations[1].Value)
				assert.Equal(t, "yet_another_code", translations[2].Code)
				assert.Equal(t, "YetAnotherValue_fr", translations[2].Value)
			})

			t.Run("Given 3 codes in 2 languages When retrieving with no language Then 3 codes have been returned", func(t *testing.T) {
				translations := cmsService.GetTranslations("")
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
		})
	})

	t.Run("Product Testing", func(t *testing.T) {
		productRepository := commonRepository.NewCrudRepository[product.ProductModel](bi.Db())
		categoryRepository := commonRepository.NewCrudRepository[category.CategoryModel](bi.Db())

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
	})
}
