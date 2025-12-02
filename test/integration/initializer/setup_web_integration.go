package initializer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danyel/ecommerce/cmd/config"
	"github.com/danyel/ecommerce/cmd/router"
	"github.com/danyel/ecommerce/internal/management"
	productmanagement "github.com/danyel/ecommerce/internal/product-management"
	shoppingbasket "github.com/danyel/ecommerce/internal/shopping-basket"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type WebIntegration struct {
	db *gorm.DB
	s  *httptest.Server
	t  *testing.T
	r  *http.Response
}

func (wi *WebIntegration) ProductManagementPostProducts(b *productmanagement.CreateProduct) *WebIntegration {
	return wi.Post(wi.forUrl("/api/product-management/v1/products"), b)
}

func (wi *WebIntegration) GetTranslations(l string) *WebIntegration {
	u := "/api/cms/v1/translations"
	if l != "" {
		u += fmt.Sprintf("?language=%s", l)
	}
	return wi.Get(wi.forUrl(u))
}

func (wi *WebIntegration) ManagementPostTranslations(b *management.CreateCms) *WebIntegration {
	return wi.Post(wi.forUrl("/api/management/v1/translations"), b)
}

func (wi *WebIntegration) ProductManagementGetProducts() *WebIntegration {
	return wi.Get(wi.forUrl("/api/product-management/v1/products"))
}

func (wi *WebIntegration) ShoppingBasketCreate() *WebIntegration {
	return wi.Post(wi.forUrl("/api/shopping-basket/v1/shopping-baskets"), nil)
}

func (wi *WebIntegration) ShoppingBasketAddItem(id string, a shoppingbasket.UpdateShoppingBasketItem) *WebIntegration {
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

func (wi *WebIntegration) Db() *gorm.DB {
	return wi.db
}

func (wi *WebIntegration) Get(url string) *WebIntegration {
	var err error
	wi.r, err = http.Get(url)
	assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) Post(url string, body any) *WebIntegration {
	b, _ := json.Marshal(body)
	var err error
	wi.r, err = http.Post(url, "application/json", bytes.NewBuffer(b))
	assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) GetResponseBody(b any) *WebIntegration {
	err := json.NewDecoder(wi.r.Body).Decode(&b)
	assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) AssertStatusCreated() *WebIntegration {
	return wi.Equal(http.StatusCreated, wi.r.StatusCode)
}

func (wi *WebIntegration) IsNotNil(b any) *WebIntegration {
	assert.NotNil(wi.t, b)
	return wi
}

func (wi *WebIntegration) AssertStatusOk() *WebIntegration {
	return wi.Equal(http.StatusOK, wi.r.StatusCode)
}

func (wi *WebIntegration) Equal(expected, actual interface{}, msgAndArgs ...interface{}) *WebIntegration {
	assert.Equal(wi.t, expected, actual, msgAndArgs...)
	return wi
}

func (wi *WebIntegration) AssertBadRequest() *WebIntegration {
	return wi.Equal(http.StatusBadRequest, wi.r.StatusCode)
}

func SetupWebIntegration(t *testing.T) *WebIntegration {
	t.Helper()
	bi := NewBackendInitializer()
	bi.Run()
	db := bi.Db()
	sc := config.NewServerConfiguration()
	ad := router.ApiDefinition{
		SC: &sc,
		DB: db,
	}

	ts := httptest.NewServer(ad.ConfigRouter())

	t.Cleanup(func() {
		ts.Close()
	})

	return &WebIntegration{
		db: db,
		s:  ts,
		t:  t,
	}
}
