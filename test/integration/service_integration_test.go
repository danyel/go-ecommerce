package integration

import (
	Testing "testing"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Product "github.com/danyel/ecommerce/internal/product"
	Initializers "github.com/danyel/ecommerce/test/integration/initializer"
	TestUtils "github.com/danyel/ecommerce/test/testutils"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	Assert "github.com/stretchr/testify/assert"
)

func TestServiceIntegration(unitTest *Testing.T) {
	TestUtils.PreInitTest()
	Logger.Log.Info("Starting Service Integration Test Cases")
	backendInitializer := Initializers.NewBackendInitializer()
	backendInitializer.TestContainers(unitTest)
	backendInitializer.Run()

	unitTest.Run("Category Testing", func(unitTest *Testing.T) {
		categoryRepository := Repository.NewCrudRepository[Category.CategoryModel](backendInitializer.DatabaseConnection())

		unitTest.Run("Get Categories Before Creation", func(unitTest *Testing.T) {
			categories := categoryRepository.FindAll(Repository.SearchCriteria{})
			Assert.Equal(unitTest, 0, len(categories))
		})

		unitTest.Run("Create Category", func(unitTest *Testing.T) {
			category := Category.CategoryModel{Name: "test", Children: []*Category.CategoryModel{}}
			err := categoryRepository.Create(&category)
			Assert.Nil(unitTest, err)
		})

		unitTest.Run("Get Categories After Creation", func(unitTest *Testing.T) {
			categories := categoryRepository.FindAll(Repository.SearchCriteria{})
			Assert.Equal(unitTest, 1, len(categories))
		})
	})

	unitTest.Run("Cms Testing", func(unitTest *Testing.T) {
		cmsRepository := Repository.NewCrudRepository[CMS.CmsModel](backendInitializer.DatabaseConnection())
		database := Database(cmsRepository)
		cmsService := CMS.NewService(cmsRepository)

		unitTest.Run("Create translation", func(unitTest *Testing.T) {
			cms := cmsRepository.Create(&CMS.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"})
			Assert.Nil(unitTest, cms)
			cms = cmsRepository.Create(&CMS.CmsModel{Code: "code", Language: "nl_fr", Value: "Value_fr"})
			Assert.Nil(unitTest, cms)
		})

		unitTest.Run("Delete translation", func(unitTest *Testing.T) {
			cms := &CMS.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"}
			err := cmsRepository.Create(cms)
			Assert.Nil(unitTest, err)
			err = cmsRepository.Delete(cms.ID)
			Assert.Nil(unitTest, err)
		})

		unitTest.Run("Get translations", func(unitTest *Testing.T) {
			database.Insert(&CMS.CmsModel{Code: "another_code", Language: "nl_be", Value: "AnotherValue_nl"})
			database.Insert(&CMS.CmsModel{Code: "another_code", Language: "nl_fr", Value: "AnotherValue_fr"})
			database.Insert(&CMS.CmsModel{Code: "yet_another_code", Language: "nl_be", Value: "YetAnotherValue_nl"})
			database.Insert(&CMS.CmsModel{Code: "yet_another_code", Language: "nl_fr", Value: "YetAnotherValue_fr"})

			unitTest.Run("Given 3 codes in 2 languages When retrieving for dutch Then 3 codes have been returned", func(unitTest *Testing.T) {
				translations := cmsService.GetTranslations("nl_be")
				Assert.Equal(unitTest, 3, len(translations))
				Assert.Equal(unitTest, "code", translations[0].Code)
				Assert.Equal(unitTest, "Value_nl", translations[0].Value)
				Assert.Equal(unitTest, "another_code", translations[1].Code)
				Assert.Equal(unitTest, "AnotherValue_nl", translations[1].Value)
				Assert.Equal(unitTest, "yet_another_code", translations[2].Code)
				Assert.Equal(unitTest, "YetAnotherValue_nl", translations[2].Value)
			})

			unitTest.Run("Given 3 codes in 2 languages When retrieving for french Then 3 codes have been returned", func(unitTest *Testing.T) {
				translations := cmsService.GetTranslations("nl_fr")
				Assert.Equal(unitTest, 3, len(translations))
				Assert.Equal(unitTest, "code", translations[0].Code)
				Assert.Equal(unitTest, "Value_fr", translations[0].Value)
				Assert.Equal(unitTest, "another_code", translations[1].Code)
				Assert.Equal(unitTest, "AnotherValue_fr", translations[1].Value)
				Assert.Equal(unitTest, "yet_another_code", translations[2].Code)
				Assert.Equal(unitTest, "YetAnotherValue_fr", translations[2].Value)
			})

			unitTest.Run("Given 3 codes in 2 languages When retrieving with no language Then 3 codes have been returned", func(unitTest *Testing.T) {
				translations := cmsService.GetTranslations("")
				Assert.Equal(unitTest, 6, len(translations))
				Assert.Equal(unitTest, "code", translations[0].Code)
				Assert.Equal(unitTest, "Value_nl", translations[0].Value)
				Assert.Equal(unitTest, "code", translations[1].Code)
				Assert.Equal(unitTest, "Value_fr", translations[1].Value)
				Assert.Equal(unitTest, "another_code", translations[2].Code)
				Assert.Equal(unitTest, "AnotherValue_nl", translations[2].Value)
				Assert.Equal(unitTest, "another_code", translations[3].Code)
				Assert.Equal(unitTest, "AnotherValue_fr", translations[3].Value)
				Assert.Equal(unitTest, "yet_another_code", translations[4].Code)
				Assert.Equal(unitTest, "YetAnotherValue_nl", translations[4].Value)
				Assert.Equal(unitTest, "yet_another_code", translations[5].Code)
				Assert.Equal(unitTest, "YetAnotherValue_fr", translations[5].Value)
			})
		})
	})

	unitTest.Run("Product Testing", func(unitTest *Testing.T) {
		productRepository := Repository.NewCrudRepository[Product.ProductModel](backendInitializer.DatabaseConnection())
		categoryRepository := Repository.NewCrudRepository[Category.CategoryModel](backendInitializer.DatabaseConnection())

		unitTest.Run("Create Product ", func(unitTest *Testing.T) {
			categoryModel := Category.CategoryModel{Name: "test", Children: []*Category.CategoryModel{}}
			err := categoryRepository.Create(&categoryModel)
			Assert.Nil(unitTest, err)
			product := Product.ProductModel{
				Brand:       "Brand",
				Name:        "Name",
				Description: "Description",
				Code:        "Code",
				Price:       1000,
				CategoryID:  categoryModel.ID,
				ImageURL:    "ImageUrl",
				Stock:       1,
			}

			err = productRepository.Create(&product)
			Assert.Nil(unitTest, err)
		})
	})
}
