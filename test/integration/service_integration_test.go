package integration

import (
	Testing "testing"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Product "github.com/danyel/ecommerce/internal/product"
	Initializers "github.com/danyel/ecommerce/test/integration/initializer"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	Assert "github.com/stretchr/testify/assert"
)

func TestServiceIntegration(t *Testing.T) {
	bi := Initializers.NewBackendInitializer()
	bi.TestContainers(t)
	bi.Run()

	t.Run("Category Testing", func(t *Testing.T) {
		categoryRepository := CommonRepository.NewCrudRepository[Category.CategoryModel](bi.Db())

		t.Run("Get Categories Before Creation", func(t *Testing.T) {
			findAll := categoryRepository.FindAll(CommonRepository.SearchCriteria{})
			Assert.Equal(t, 0, len(findAll))
		})

		t.Run("Create Category", func(t *Testing.T) {
			c := Category.CategoryModel{Name: "test", Children: []*Category.CategoryModel{}}
			e := categoryRepository.Create(&c)
			Assert.Nil(t, e)
		})

		t.Run("Get Categories After Creation", func(t *Testing.T) {
			findAll := categoryRepository.FindAll(CommonRepository.SearchCriteria{})
			Assert.Equal(t, 1, len(findAll))
		})
	})

	t.Run("Cms Testing", func(t *Testing.T) {
		cmsRepository := CommonRepository.NewCrudRepository[CMS.CmsModel](bi.Db())
		f := Database[CMS.CmsModel](cmsRepository)
		cmsService := CMS.NewCmsService(bi.Db())

		t.Run("Create translation", func(t *Testing.T) {
			e := cmsRepository.Create(&CMS.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"})
			Assert.Nil(t, e)
			e = cmsRepository.Create(&CMS.CmsModel{Code: "code", Language: "nl_fr", Value: "Value_fr"})
			Assert.Nil(t, e)
		})

		t.Run("Delete translation", func(t *Testing.T) {
			p := &CMS.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"}
			e := cmsRepository.Create(p)
			Assert.Nil(t, e)
			err := cmsRepository.Delete(p.ID)
			Assert.Nil(t, err)
		})

		t.Run("Get translations", func(t *Testing.T) {
			f.Insert(&CMS.CmsModel{Code: "another_code", Language: "nl_be", Value: "AnotherValue_nl"})
			f.Insert(&CMS.CmsModel{Code: "another_code", Language: "nl_fr", Value: "AnotherValue_fr"})
			f.Insert(&CMS.CmsModel{Code: "yet_another_code", Language: "nl_be", Value: "YetAnotherValue_nl"})
			f.Insert(&CMS.CmsModel{Code: "yet_another_code", Language: "nl_fr", Value: "YetAnotherValue_fr"})

			t.Run("Given 3 codes in 2 languages When retrieving for dutch Then 3 codes have been returned", func(t *Testing.T) {
				translations := cmsService.GetTranslations("nl_be")
				Assert.Equal(t, 3, len(translations))
				Assert.Equal(t, "code", translations[0].Code)
				Assert.Equal(t, "Value_nl", translations[0].Value)
				Assert.Equal(t, "another_code", translations[1].Code)
				Assert.Equal(t, "AnotherValue_nl", translations[1].Value)
				Assert.Equal(t, "yet_another_code", translations[2].Code)
				Assert.Equal(t, "YetAnotherValue_nl", translations[2].Value)
			})

			t.Run("Given 3 codes in 2 languages When retrieving for french Then 3 codes have been returned", func(t *Testing.T) {
				translations := cmsService.GetTranslations("nl_fr")
				Assert.Equal(t, 3, len(translations))
				Assert.Equal(t, "code", translations[0].Code)
				Assert.Equal(t, "Value_fr", translations[0].Value)
				Assert.Equal(t, "another_code", translations[1].Code)
				Assert.Equal(t, "AnotherValue_fr", translations[1].Value)
				Assert.Equal(t, "yet_another_code", translations[2].Code)
				Assert.Equal(t, "YetAnotherValue_fr", translations[2].Value)
			})

			t.Run("Given 3 codes in 2 languages When retrieving with no language Then 3 codes have been returned", func(t *Testing.T) {
				translations := cmsService.GetTranslations("")
				Assert.Equal(t, 6, len(translations))
				Assert.Equal(t, "code", translations[0].Code)
				Assert.Equal(t, "Value_nl", translations[0].Value)
				Assert.Equal(t, "code", translations[1].Code)
				Assert.Equal(t, "Value_fr", translations[1].Value)
				Assert.Equal(t, "another_code", translations[2].Code)
				Assert.Equal(t, "AnotherValue_nl", translations[2].Value)
				Assert.Equal(t, "another_code", translations[3].Code)
				Assert.Equal(t, "AnotherValue_fr", translations[3].Value)
				Assert.Equal(t, "yet_another_code", translations[4].Code)
				Assert.Equal(t, "YetAnotherValue_nl", translations[4].Value)
				Assert.Equal(t, "yet_another_code", translations[5].Code)
				Assert.Equal(t, "YetAnotherValue_fr", translations[5].Value)
			})
		})
	})

	t.Run("Product Testing", func(t *Testing.T) {
		productRepository := CommonRepository.NewCrudRepository[Product.ProductModel](bi.Db())
		categoryRepository := CommonRepository.NewCrudRepository[Category.CategoryModel](bi.Db())

		t.Run("Create Product ", func(t *Testing.T) {
			c := Category.CategoryModel{Name: "test", Children: []*Category.CategoryModel{}}
			e := categoryRepository.Create(&c)
			Assert.Nil(t, e)
			p := Product.ProductModel{
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
			Assert.Nil(t, err)
		})
	})
}
