package integration

import (
	Log "log"
	Testing "testing"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	Shoppingbasket "github.com/danyel/ecommerce/internal/shoppingbasket"
	Initializers "github.com/danyel/ecommerce/test/integration/initializer"
)

func TestWebHandler(unitTest *Testing.T) {
	Log.Print("Starting Web Handler Test Cases")
	webIntegration := Initializers.SetupWebIntegration(unitTest)
	categoryRepository := Database(Repository.NewCrudRepository[Category.CategoryModel](webIntegration.DatabaseConnection()))
	cmsRepository := Database(Repository.NewCrudRepository[CMS.CmsModel](webIntegration.DatabaseConnection()))
	productRepository := Database(Repository.NewCrudRepository[Product.ProductModel](webIntegration.DatabaseConnection()))

	categoryModel := &Category.CategoryModel{Name: "GPU", Children: []*Category.CategoryModel{}}
	categoryRepository.Insert(categoryModel)
	cmsRepository.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_NAME",
		Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart",
		Language: "nl_BE",
	})
	cmsRepository.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_DESCRIPTION",
		Value:    "De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.",
		Language: "nl_BE",
	})
	cmsRepository.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_NAME_FR",
		Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart_FR",
		Language: "nl_FR",
	})
	cmsRepository.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_DESCRIPTION_FR",
		Value:    "De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.FR",
		Language: "nl_FR",
	})
	productModel := &Product.ProductModel{
		Brand:       "ASUS",
		Name:        "90YV0L71_M0NA00_NAME",
		Description: "90YV0L71_M0NA00_DESCRIPTION",
		Code:        "90YV0L71-M0NA00",
		Price:       669,
		CategoryID:  categoryModel.ID,
		ImageURL:    "https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg",
	}
	productRepository.Insert(productModel)

	unitTest.Run("Product h", func(unitTest *Testing.T) {
		unitTest.Run("CreateProduct", func(unitTest *Testing.T) {
			createProduct := &Product.CreateProduct{
				Brand:       "ASUS",
				Name:        "90YV0L71_M0NA00_NAME",
				Description: "90YV0L71_M0NA00_DESCRIPTION",
				Code:        "90YV0L71-M0NA00",
				Price:       669,
				CategoryID:  Types.NewID(categoryModel.ID),
				ImageURL:    "https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg",
			}
			var ID Types.ID
			webIntegration.ProductManagementPostProducts(createProduct).
				GetResponseBody(&ID).
				AssertStatusCreated().
				IsNotNil(ID.ID)
		})
	})

	unitTest.Run("CMS h", func(unitTest *Testing.T) {
		unitTest.Run("TestCmsHandler", func(unitTest *Testing.T) {
			unitTest.Run("CmsHandler retrieve dutch", func(unitTest *Testing.T) {
				var translations []CMS.Translation
				webIntegration.GetTranslations("nl_BE").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(2, len(translations)).
					Equal("90YV0L71_M0NA00_NAME", translations[0].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", translations[0].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION", translations[1].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", translations[1].Value)
			})

			unitTest.Run("CmsHandler retrieve french", func(unitTest *Testing.T) {
				var translations []CMS.Translation
				webIntegration.GetTranslations("nl_FR").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(2, len(translations)).
					Equal("90YV0L71_M0NA00_NAME_FR", translations[0].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart_FR", translations[0].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION_FR", translations[1].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.FR", translations[1].Value)
			})

			unitTest.Run("CmsHandler retrieve all", func(unitTest *Testing.T) {
				var translations []CMS.Translation
				webIntegration.GetTranslations("").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(4, len(translations)).
					Equal("90YV0L71_M0NA00_NAME", translations[0].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", translations[0].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION", translations[1].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", translations[1].Value).
					Equal("90YV0L71_M0NA00_NAME_FR", translations[2].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart_FR", translations[2].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION_FR", translations[3].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.FR", translations[3].Value)
			})

			unitTest.Run("CmsHandler retrieve none because of invalid language", func(unitTest *Testing.T) {
				var translations []CMS.Translation
				webIntegration.GetTranslations("nl_de").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(0, len(translations))
			})
		})
	})

	unitTest.Run("Management h", func(unitTest *Testing.T) {
		unitTest.Run("ManagementHandler: Create a new translation but it already exist so return 400", func(unitTest *Testing.T) {
			createCms := &Management.CreateCms{
				Code:     "90YV0L71_M0NA00_NAME",
				Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart",
				Language: "nl_BE",
			}
			webIntegration.ManagementPostTranslations(createCms).
				AssertBadRequest()
		})

		unitTest.Run("ManagementHandler: Create a new translation and 201 is return", func(unitTest *Testing.T) {
			createCms := &Management.CreateCms{
				Code:     "unknown",
				Language: "nl_fr",
				Value:    "Value_fr",
			}
			var i Management.CmsID
			webIntegration.ManagementPostTranslations(createCms).
				GetResponseBody(&i).
				IsNotNil(i.ID).
				AssertStatusCreated()
		})
	})

	unitTest.Run("Product Management h", func(unitTest *Testing.T) {
		unitTest.Run("Product Management Get Product", func(unitTest *Testing.T) {
			var product Product.Product
			webIntegration.ProductManagementGetProductByID(productModel.ID.String()).
				GetResponseBody(&product).
				AssertStatusOk().
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", product.Name).
				Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", product.Description).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", product.ImageURL).
				Equal("GPU", product.Category.Name).
				Equal(Types.Float64(669), product.Price.Inclusive).
				Equal("90YV0L71-M0NA00", product.Code).
				Equal("ASUS", product.Brand).
				Equal(productModel.ID, product.ID.ID)
		})

		unitTest.Run("Product Management Get Products", func(unitTest *Testing.T) {
			var products []Product.Product
			webIntegration.ProductManagementGetProducts().
				GetResponseBody(&products).
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", products[0].Name).
				Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", products[0].Description).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", products[0].ImageURL).
				Equal("GPU", products[0].Category.Name).
				Equal(Types.Float64(669), products[0].Price.Inclusive).
				Equal("90YV0L71-M0NA00", products[0].Code).
				Equal("ASUS", products[0].Brand).
				Equal(productModel.ID, products[0].ID.ID).
				AssertStatusOk()
		})
	})

	unitTest.Run("Shopping Basket h", func(unitTest *Testing.T) {
		var ID Types.ID

		unitTest.Run("Create Shopping Basket", func(unitTest *Testing.T) {
			webIntegration.ShoppingBasketCreate().
				GetResponseBody(&ID).
				AssertStatusCreated().
				IsNotNil(ID)
		})

		unitTest.Run("Add Item To Shopping Basket", func(unitTest *Testing.T) {
			updateShoppingBasketItem := Shoppingbasket.UpdateShoppingBasketItem{
				ProductID: Types.NewID(productModel.ID),
			}
			webIntegration.ShoppingBasketAddItem(ID.ID.String(), updateShoppingBasketItem).
				AssertStatusOk()
		})

		unitTest.Run("Get Shopping Basket", func(unitTest *Testing.T) {
			var shoppingBasket Shoppingbasket.ShoppingBasket
			webIntegration.GetShoppingBasket(ID.ID.String()).
				GetResponseBody(&shoppingBasket).
				AssertStatusOk().
				Equal(ID, shoppingBasket.ID).
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", shoppingBasket.Items[0].Name).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", shoppingBasket.Items[0].ImageURL).
				Equal(Types.Float64(669.0), shoppingBasket.Items[0].BasePrice.Inclusive)
		})
	})
}
