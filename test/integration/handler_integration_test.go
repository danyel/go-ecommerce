package integration

import (
	"testing"

	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/internal/management"
	"github.com/dnoulet/ecommerce/internal/product"
	productmanagement "github.com/dnoulet/ecommerce/internal/product-management"
	"github.com/dnoulet/ecommerce/test/integration/initializer"
)

func TestHandler(t *testing.T) {
	wi := initializer.SetupWebIntegration(t)

	t.Run("Product Handler", func(t *testing.T) {
		d := Database(commonRepository.NewCrudRepository[category.CategoryModel](wi.Db()))

		t.Run("CreateProduct", func(t *testing.T) {
			c := &category.CategoryModel{Name: "test", Children: []*category.CategoryModel{}}
			d.Insert(c)
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
		f := Database(commonRepository.NewCrudRepository[cms.CmsModel](wi.Db()))
		f.Insert(&cms.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"})
		f.Insert(&cms.CmsModel{Code: "code", Language: "nl_fr", Value: "Value_fr"})
		f.Insert(&cms.CmsModel{Code: "another_code", Language: "nl_be", Value: "AnotherValue_nl"})
		f.Insert(&cms.CmsModel{Code: "another_code", Language: "nl_fr", Value: "AnotherValue_fr"})
		f.Insert(&cms.CmsModel{Code: "yet_another_code", Language: "nl_be", Value: "YetAnotherValue_nl"})
		f.Insert(&cms.CmsModel{Code: "yet_another_code", Language: "nl_fr", Value: "YetAnotherValue_fr"})

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

	t.Run("Product Management Handler", func(t *testing.T) {
		c := Database(commonRepository.NewCrudRepository[cms.CmsModel](wi.Db()))
		p := Database(commonRepository.NewCrudRepository[product.ProductModel](wi.Db()))
		d := Database(commonRepository.NewCrudRepository[category.CategoryModel](wi.Db()))
		c.Insert(&cms.CmsModel{
			Code:     "90YV0L71_M0NA00_NAME",
			Value:    "MSI Prime Radeon RX 9070 XT 16GB OC Videokaart",
			Language: "nl_BE",
		})
		c.Insert(&cms.CmsModel{
			Code:     "90YV0L71_M0NA00_DESCRIPTION",
			Value:    "De ASUS Prime Radeon RX 9070 XT Gaming OC 16GB Videokaart is een krachtige AMD-kaart die is uitgerust met 16 GB GDDR6-videogeheugen en een GPU-kloksnelheid van tot wel 3030 MHz. Met 4096 stream processors biedt deze videokaart uitstekende prestaties voor zowel gaming als professionele toepassingen. De ASUS Prime-serie is ontworpen voor gamers en enthousiastelingen die op zoek zijn naar een betrouwbare en geavanceerde grafische oplossing.",
			Language: "nl_BE",
		})
		cm := &category.CategoryModel{Name: "GPU", Children: []*category.CategoryModel{}}
		d.Insert(cm)
		pm := &product.ProductModel{
			Brand:       "ASUS",
			Name:        "90YV0L71_M0NA00_NAME",
			Description: "90YV0L71_M0NA00_DESCRIPTION",
			Code:        "90YV0L71-M0NA00",
			Price:       669,
			CategoryId:  cm.ID,
			ImageUrl:    "https://www.megekko.nl/productimg/1699548/nw/1_ASUS-Prime-Radeon-RX-9070-XT-16GB-OC-Videokaart.jpg",
		}
		p.Insert(pm)

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
}
