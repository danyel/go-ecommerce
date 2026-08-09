package initializer

import (
	Bytes "bytes"
	JSON "encoding/json"
	Fmt "fmt"
	Log "log"
	Http "net/http"
	HttpTest "net/http/httptest"
	Testing "testing"

	Broker "github.com/danyel/ecommerce/cmd/broker"
	Configuration "github.com/danyel/ecommerce/cmd/config"
	ApplicationRouter "github.com/danyel/ecommerce/cmd/router"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shopping-basket"
	Assert "github.com/stretchr/testify/assert"
	Database "gorm.io/gorm"
)

type WebIntegration struct {
	db *Database.DB
	s  *HttpTest.Server
	t  *Testing.T
	r  *Http.Response
}

func (wi *WebIntegration) ProductManagementPostProducts(b *Product.CreateProduct) *WebIntegration {
	return wi.Post(wi.forUrl("/api/product-management/v1/products"), b)
}

func (wi *WebIntegration) GetTranslations(l string) *WebIntegration {
	u := "/api/cms/v1/translations"
	if l != "" {
		u += Fmt.Sprintf("?language=%s", l)
	}
	return wi.Get(wi.forUrl(u))
}

func (wi *WebIntegration) ManagementPostTranslations(b *Management.CreateCms) *WebIntegration {
	return wi.Post(wi.forUrl("/api/management/v1/translations"), b)
}

func (wi *WebIntegration) ProductManagementGetProducts() *WebIntegration {
	return wi.Get(wi.forUrl("/api/product-management/v1/products"))
}

func (wi *WebIntegration) ShoppingBasketCreate() *WebIntegration {
	return wi.Post(wi.forUrl("/api/shopping-basket/v1/shopping-baskets"), nil)
}

func (wi *WebIntegration) ShoppingBasketAddItem(id string, a ShoppingBasket.UpdateShoppingBasketItem) *WebIntegration {
	return wi.Post(wi.forUrl("/api/shopping-basket/v1/shopping-baskets/"+id), a)
}

func (wi *WebIntegration) GetShoppingBasket(id string) *WebIntegration {
	return wi.Get(wi.forUrl("/api/shopping-basket/v1/shopping-baskets/" + id))
}

func (wi *WebIntegration) ProductManagementGetProductById(i string) *WebIntegration {
	return wi.Get(wi.forUrl("/api/product-management/v1/products/" + i))
}

func (wi *WebIntegration) forUrl(url string) string {
	return wi.s.URL + url
}

func (wi *WebIntegration) Db() *Database.DB {
	return wi.db
}

func (wi *WebIntegration) Get(url string) *WebIntegration {
	var err error
	wi.r, err = Http.Get(url)
	Assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) Post(url string, body any) *WebIntegration {
	b, _ := JSON.Marshal(body)
	var err error
	wi.r, err = Http.Post(url, "application/json", Bytes.NewBuffer(b))
	Assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) GetResponseBody(b any) *WebIntegration {
	err := JSON.NewDecoder(wi.r.Body).Decode(&b)
	Assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) AssertStatusCreated() *WebIntegration {
	return wi.Equal(Http.StatusCreated, wi.r.StatusCode)
}

func (wi *WebIntegration) IsNotNil(b any) *WebIntegration {
	Assert.NotNil(wi.t, b)
	return wi
}

func (wi *WebIntegration) AssertStatusOk() *WebIntegration {
	return wi.Equal(Http.StatusOK, wi.r.StatusCode)
}

func (wi *WebIntegration) Equal(expected, actual any, msgAndArgs ...any) *WebIntegration {
	Assert.Equal(wi.t, expected, actual, msgAndArgs...)
	return wi
}

func (wi *WebIntegration) AssertBadRequest() *WebIntegration {
	return wi.Equal(Http.StatusBadRequest, wi.r.StatusCode)
}

func SetupWebIntegration(t *Testing.T) *WebIntegration {
	t.Helper()
	bi := NewBackendInitializer()
	bi.TestContainers(t)
	bi.Run()
	db := bi.Db()
	sc := Configuration.NewServerConfiguration()
	newBroker := Broker.NewBroker()
	err := newBroker.CreateConnection(bi.BrokerConfiguration)
	if err != nil {
		t.Fatal(err)
	}
	ShoppingBasket.RegisterShoppingBasketEvents(ShoppingBasket.NewService(db, newBroker), newBroker)
	if err := newBroker.Start(); err != nil {
		Log.Println(err.Error())
	}
	ad := ApplicationRouter.ApiDefinition{
		SC:             &sc,
		DB:             db,
		EventPublisher: newBroker,
	}

	ts := HttpTest.NewServer(ad.ConfigRouter())

	t.Cleanup(func() {
		ts.Close()
	})

	return &WebIntegration{
		db: db,
		s:  ts,
		t:  t,
	}
}
