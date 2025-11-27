package integration

import (
	"testing"

	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"github.com/danyel/ecommerce/internal/management"
	"github.com/danyel/ecommerce/internal/product"
	productmanagement "github.com/danyel/ecommerce/internal/product-management"
	shoppingbasket "github.com/danyel/ecommerce/internal/shopping-basket"
	"github.com/danyel/ecommerce/test/integration/initializer"
)

func TestHandler(t *testing.T) {
	wi := initializer.SetupWebIntegration(t)
	categoryRepo := Database(commonRepository.NewCrudRepository[category.CategoryModel](wi.Db()))
	cmsRepo := Database(commonRepository.NewCrudRepository[cms.CmsModel](wi.Db()))
	productRepo := Database(commonRepository.NewCrudRepository[product.ProductModel](wi.Db()))

	cm := &category.CategoryModel{Name: "GPU", Children: []*category.CategoryModel{}}
	categoryRepo.Insert(cm)
	cmsRepo.Insert(&cms.CmsModel{
		Code:     "90YV0L71_M0NA00_NAME",
		Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart",
		Language: "nl_BE",
	})
	cmsRepo.Insert(&cms.CmsModel{
		Code:     "90YV0L71_M0NA00_DESCRIPTION",
		Value:    "De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.",
		Language: "nl_BE",
	})
	cmsRepo.Insert(&cms.CmsModel{
		Code:     "90YV0L71_M0NA00_NAME_FR",
		Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart_FR",
		Language: "nl_FR",
	})
	cmsRepo.Insert(&cms.CmsModel{
		Code:     "90YV0L71_M0NA00_DESCRIPTION_FR",
		Value:    "De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.FR",
		Language: "nl_FR",
	})
	pm := &product.ProductModel{
		Brand:       "ASUS",
		Name:        "90YV0L71_M0NA00_NAME",
		Description: "90YV0L71_M0NA00_DESCRIPTION",
		Code:        "90YV0L71-M0NA00",
		Price:       669,
		CategoryId:  cm.ID,
		ImageUrl:    "https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg",
	}
	productRepo.Insert(pm)

	t.Run("Product Handler", func(t *testing.T) {
		t.Run("CreateProduct", func(t *testing.T) {
			b := &productmanagement.CreateProduct{
				Brand:       "ASUS",
				Name:        "90YV0L71_M0NA00_NAME",
				Description: "90YV0L71_M0NA00_DESCRIPTION",
				Code:        "90YV0L71-M0NA00",
				Price:       669,
				CategoryId:  cm.ID,
				ImageUrl:    "https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg",
			}
			var productId productmanagement.ProductId
			wi.ProductManagementPostProducts(b).
				GetResponseBody(&productId).
				AssertStatusCreated().
				IsNotNil(productId.ID)
		})
	})

	t.Run("CMS Handler", func(t *testing.T) {
		t.Run("TestCmsHandler", func(t *testing.T) {
			t.Run("CmsHandler retrieve dutch", func(t *testing.T) {
				var translations []cms.Translation
				wi.GetTranslations("nl_BE").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(2, len(translations)).
					Equal("90YV0L71_M0NA00_NAME", translations[0].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", translations[0].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION", translations[1].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", translations[1].Value)
			})

			t.Run("CmsHandler retrieve french", func(t *testing.T) {
				var translations []cms.Translation
				wi.GetTranslations("nl_FR").
					GetResponseBody(&translations).
					AssertStatusOk().
					Equal(2, len(translations)).
					Equal("90YV0L71_M0NA00_NAME_FR", translations[0].Code).
					Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart_FR", translations[0].Value).
					Equal("90YV0L71_M0NA00_DESCRIPTION_FR", translations[1].Code).
					Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.FR", translations[1].Value)
			})

			t.Run("CmsHandler retrieve all", func(t *testing.T) {
				var translations []cms.Translation
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
				Code:     "90YV0L71_M0NA00_NAME",
				Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart",
				Language: "nl_BE",
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

	t.Run("Product Management Handler", func(t *testing.T) {
		t.Run("Product Management Get Product", func(t *testing.T) {
			var ps productmanagement.Product
			wi.ProductManagementGetProductById(pm.ID.String()).
				GetResponseBody(&ps).
				AssertStatusOk().
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", ps.Name).
				Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", ps.Description).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", ps.ImageUrl).
				Equal("GPU", ps.Category.Name).
				Equal(669, ps.Price).
				Equal("90YV0L71-M0NA00", ps.Code).
				Equal("ASUS", ps.Brand).
				Equal(pm.ID, ps.ID)
		})

		t.Run("Product Management Get Products", func(t *testing.T) {
			var ps []productmanagement.Product
			wi.ProductManagementGetProducts().
				GetResponseBody(&ps).
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", ps[0].Name).
				Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", ps[0].Description).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", ps[0].ImageUrl).
				Equal("GPU", ps[0].Category.Name).
				Equal(669, ps[0].Price).
				Equal("90YV0L71-M0NA00", ps[0].Code).
				Equal("ASUS", ps[0].Brand).
				Equal(pm.ID, ps[0].ID).
				AssertStatusOk()
		})
	})

	t.Run("Shopping Basket Handler", func(t *testing.T) {
		var shoppingId shoppingbasket.ShoppingId

		t.Run("Create Shopping Basket", func(t *testing.T) {
			wi.ShoppingBasketCreate().
				GetResponseBody(&shoppingId).
				AssertStatusCreated().
				IsNotNil(shoppingId)
		})

		t.Run("Add Item To Shopping Basket", func(t *testing.T) {
			a := shoppingbasket.AddItem{
				ProductId: pm.ID,
			}
			wi.ShoppingBasketAddItem(shoppingId.Id.String(), a).
				AssertStatusOk()
		})

		t.Run("Get Shopping Basket", func(t *testing.T) {
			var s shoppingbasket.ShoppingBasket
			wi.GetShoppingBasket(shoppingId.Id.String()).
				GetResponseBody(&s).
				AssertStatusOk().
				Equal(shoppingId.Id, s.Id).
				Equal("MSI Prime Radeon RX 9070 XT 16GB OC Videokaart", s.Items[0].Name).
				Equal("De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.", s.Items[0].Description).
				Equal("https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg", s.Items[0].ImageUrl).
				Equal("GPU", s.Items[0].Category.Name).
				Equal(669, s.Items[0].Price).
				Equal("90YV0L71-M0NA00", s.Items[0].Code).
				Equal("ASUS", s.Items[0].Brand)
		})
	})
}
