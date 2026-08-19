package integration

import (
	Testing "testing"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	Shoppingbasket "github.com/danyel/ecommerce/internal/shoppingbasket"
	Initializers "github.com/danyel/ecommerce/test/integration/initializer"
)

func TestHandler(t *Testing.T) {
	wi := Initializers.SetupWebIntegration(t)
	categoryRepo := Database(CommonRepository.NewCrudRepository[Category.CategoryModel](wi.DB()))
	cmsRepo := Database(CommonRepository.NewCrudRepository[CMS.CmsModel](wi.DB()))
	productRepo := Database(CommonRepository.NewCrudRepository[Product.ProductModel](wi.DB()))

	cm := &Category.CategoryModel{Name: "GPU", Children: []*Category.CategoryModel{}}
	categoryRepo.Insert(cm)
	cmsRepo.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_NAME",
		Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart",
		Language: "nl_BE",
	})
	cmsRepo.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_DESCRIPTION",
		Value:    "De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.",
		Language: "nl_BE",
	})
	cmsRepo.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_NAME_FR",
		Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart_FR",
		Language: "nl_FR",
	})
	cmsRepo.Insert(&CMS.CmsModel{
		Code:     "90YV0L71_M0NA00_DESCRIPTION_FR",
		Value:    "De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.FR",
		Language: "nl_FR",
	})
	pm := &Product.ProductModel{
		Brand:       "ASUS",
		Name:        "90YV0L71_M0NA00_NAME",
		Description: "90YV0L71_M0NA00_DESCRIPTION",
		Code:        "90YV0L71-M0NA00",
		Price:       669,
		CategoryID:  cm.ID,
		ImageURL:    "https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg",
	}
	productRepo.Insert(pm)

	t.Run("Product h", func(t *Testing.T) {
		t.Run("CreateProduct", func(t *Testing.T) {
			b := &Product.CreateProduct{
				Brand:       "ASUS",
				Name:        "90YV0L71_M0NA00_NAME",
				Description: "90YV0L71_M0NA00_DESCRIPTION",
				Code:        "90YV0L71-M0NA00",
				Price:       669,
				CategoryID:  Types.NewID(cm.ID),
				ImageURL:    "https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg",
			}
			var productID Types.ID
			wi.ProductManagementPostProducts(b).
				GetResponseBody(&productID).
				AssertStatusCreated().
				IsNotNil(productID.ID)
		})
	})

	t.Run("CMS h", func(t *Testing.T) {
		t.Run("TestCmsHandler", func(t *Testing.T) {
			t.Run("CmsHandler retrieve dutch", func(t *Testing.T) {
				var translations []CMS.Translation
				wi.GetTranslations("nl_BE").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(2, len(translations)).
					Equal("90YV0L71_M0NA00_NAME", translations[0].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", translations[0].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION", translations[1].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", translations[1].Value)
			})

			t.Run("CmsHandler retrieve french", func(t *Testing.T) {
				var translations []CMS.Translation
				wi.GetTranslations("nl_FR").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(2, len(translations)).
					Equal("90YV0L71_M0NA00_NAME_FR", translations[0].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart_FR", translations[0].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION_FR", translations[1].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.FR", translations[1].Value)
			})

			t.Run("CmsHandler retrieve all", func(t *Testing.T) {
				var translations []CMS.Translation
				wi.GetTranslations("").
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

			t.Run("CmsHandler retrieve none because of invalid language", func(t *Testing.T) {
				var translations []CMS.Translation
				wi.GetTranslations("nl_de").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(0, len(translations))
			})
		})
	})

	t.Run("Management h", func(t *Testing.T) {
		t.Run("ManagementHandler: Create a new translation but it already exist so return 400", func(t *Testing.T) {
			b := &Management.CreateCms{
				Code:     "90YV0L71_M0NA00_NAME",
				Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart",
				Language: "nl_BE",
			}
			wi.ManagementPostTranslations(b).
				AssertBadRequest()
		})

		t.Run("ManagementHandler: Create a new translation and 201 is return", func(t *Testing.T) {
			b := &Management.CreateCms{
				Code:     "unknown",
				Language: "nl_fr",
				Value:    "Value_fr",
			}
			var i Management.CmsId
			wi.ManagementPostTranslations(b).
				GetResponseBody(&i).
				IsNotNil(i.ID).
				AssertStatusCreated()
		})
	})

	t.Run("Product Management h", func(t *Testing.T) {
		t.Run("Product Management Get Product", func(t *Testing.T) {
			var ps Product.Product
			wi.ProductManagementGetProductByID(pm.ID.String()).
				GetResponseBody(&ps).
				AssertStatusOk().
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", ps.Name).
				Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", ps.Description).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", ps.ImageURL).
				Equal("GPU", ps.Category.Name).
				Equal(Types.Float64(669), ps.Price.Inclusive).
				Equal("90YV0L71-M0NA00", ps.Code).
				Equal("ASUS", ps.Brand).
				Equal(pm.ID, ps.ID.ID)
		})

		t.Run("Product Management Get Products", func(t *Testing.T) {
			var ps []Product.Product
			wi.ProductManagementGetProducts().
				GetResponseBody(&ps).
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", ps[0].Name).
				Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", ps[0].Description).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", ps[0].ImageURL).
				Equal("GPU", ps[0].Category.Name).
				Equal(Types.Float64(669), ps[0].Price.Inclusive).
				Equal("90YV0L71-M0NA00", ps[0].Code).
				Equal("ASUS", ps[0].Brand).
				Equal(pm.ID, ps[0].ID.ID).
				AssertStatusOk()
		})
	})

	t.Run("Shopping Basket h", func(t *Testing.T) {
		var shoppingID Types.ID

		t.Run("Create Shopping Basket", func(t *Testing.T) {
			wi.ShoppingBasketCreate().
				GetResponseBody(&shoppingID).
				AssertStatusCreated().
				IsNotNil(shoppingID)
		})

		t.Run("Add Item To Shopping Basket", func(t *Testing.T) {
			a := Shoppingbasket.UpdateShoppingBasketItem{
				ProductID: Types.NewID(pm.ID),
			}
			wi.ShoppingBasketAddItem(shoppingID.ID.String(), a).
				AssertStatusOk()
		})

		t.Run("Get Shopping Basket", func(t *Testing.T) {
			var s Shoppingbasket.ShoppingBasket
			wi.GetShoppingBasket(shoppingID.ID.String()).
				GetResponseBody(&s).
				AssertStatusOk().
				Equal(shoppingID, s.ID).
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", s.Items[0].Name).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", s.Items[0].ImageURL).
				Equal(Types.Float64(669.0), s.Items[0].BasePrice.Inclusive)
		})
	})
}
