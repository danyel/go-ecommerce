package integration

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
	DB     *gorm.DB
	Server *httptest.Server
	t      *testing.T
	r      *http.Response
}

type StatusChecker struct {
}

func (ta *WebIntegration) Get(url string) *WebIntegration {
	var err error
	ta.r, err = http.Get(url)
	assert.Nil(ta.t, err)
	return ta
}

func (ta *WebIntegration) Post(url, contentType string, body any) *WebIntegration {
	b, _ := json.Marshal(body)
	var err error
	ta.r, err = http.Post(url, contentType, bytes.NewBuffer(b))
	assert.Nil(ta.t, err)
	return ta
}

func (ta *WebIntegration) GetResponseBody(b any) *WebIntegration {
	err := json.NewDecoder(ta.r.Body).Decode(&b)
	assert.Nil(ta.t, err)
	return ta
}

func (ta *WebIntegration) AssertStatusCreated() {
	assert.Equal(ta.t, http.StatusCreated, ta.r.StatusCode)
}

func (ta *WebIntegration) AssertStatusOk() {
	assert.Equal(ta.t, http.StatusOK, ta.r.StatusCode)
}

func (ta *WebIntegration) AssertBadRequest() {
	assert.Equal(ta.t, http.StatusBadRequest, ta.r.StatusCode)
}

func SetupWebIntegration(t *testing.T) *WebIntegration {
	t.Helper()
	initializer := NewBackendInitializer()
	initializer.Run()
	db := initializer.DB
	serverConfiguration := config.NewServerConfiguration()
	apiDefinition := router.ApiDefinition{
		ServerConfiguration: &serverConfiguration,
		DB:                  db,
	}

	// Create a test server
	ts := httptest.NewServer(apiDefinition.ConfigRouter())

	t.Cleanup(func() {
		ts.Close()
	})

	return &WebIntegration{
		DB:     db,
		Server: ts,
		t:      t,
	}
}
