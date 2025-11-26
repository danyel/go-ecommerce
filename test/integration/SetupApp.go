package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/dnoulet/ecommerce/cmd/config"
	"github.com/dnoulet/ecommerce/cmd/router"
	"gorm.io/gorm"
)

type TestApp struct {
	DB     *gorm.DB
	Server *httptest.Server
}

func SetupTestApp(t *testing.T) *TestApp {
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

	return &TestApp{
		DB:     db,
		Server: ts,
	}
}
