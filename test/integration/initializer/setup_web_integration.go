package initializer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dnoulet/ecommerce/cmd/config"
	"github.com/dnoulet/ecommerce/cmd/router"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type WebIntegration struct {
	db *gorm.DB
	s  *httptest.Server
	t  *testing.T
	r  *http.Response
}

func (wi *WebIntegration) ForUrl(url string) string {
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

func (wi *WebIntegration) Post(url, contentType string, body any) *WebIntegration {
	b, _ := json.Marshal(body)
	var err error
	wi.r, err = http.Post(url, contentType, bytes.NewBuffer(b))
	assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) GetResponseBody(b any) *WebIntegration {
	err := json.NewDecoder(wi.r.Body).Decode(&b)
	assert.Nil(wi.t, err)
	return wi
}

func (wi *WebIntegration) AssertStatusCreated() {
	assert.Equal(wi.t, http.StatusCreated, wi.r.StatusCode)
}

func (wi *WebIntegration) AssertStatusOk() {
	assert.Equal(wi.t, http.StatusOK, wi.r.StatusCode)
}

func (wi *WebIntegration) AssertBadRequest() {
	assert.Equal(wi.t, http.StatusBadRequest, wi.r.StatusCode)
}

func SetupWebIntegration(t *testing.T) *WebIntegration {
	t.Helper()
	bi := NewBackendInitializer()
	bi.Run()
	db := bi.Db()
	sc := config.NewServerConfiguration()
	ad := router.ApiDefinition{
		ServerConfiguration: &sc,
		DB:                  db,
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
