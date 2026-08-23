package initializer

import (
	"bytes"
	JSON "encoding/json"
	Fmt "fmt"
	"io"
	Log "log"
	Http "net/http"
	HttpTest "net/http/httptest"
	Testing "testing"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	ApplicationRouter "github.com/danyel/ecommerce/cmd/router"
	Factory "github.com/danyel/ecommerce/internal/common/factory"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shoppingbasket"
	Assert "github.com/stretchr/testify/assert"
	Database "gorm.io/gorm"
)

type WebIntegration struct {
	databaseConnection *Database.DB
	server             *HttpTest.Server
	unitTest           *Testing.T
	response           *Http.Response
	authToken          string
}

func (webIntegration *WebIntegration) WithAuth(ID string, roles []string, hexSecret string) *WebIntegration {
	useClaims := &ApplicationMiddleware.UserClaims{
		UserID: ID,
		Roles:  roles,
	}

	authToken, err := ApplicationMiddleware.EncryptClaims(useClaims, hexSecret)
	Assert.Nil(webIntegration.unitTest, err)

	webIntegration.authToken = authToken
	return webIntegration
}

func (webIntegration *WebIntegration) ProductManagementPostProducts(createProduct *Product.CreateProduct) *WebIntegration {
	return webIntegration.Post(webIntegration.forURL("/api/product-management/v1/products"), createProduct)
}

func (webIntegration *WebIntegration) GetTranslations(language string) *WebIntegration {
	baseUrl := "/api/cms/v1/translations"
	if language != "" {
		baseUrl += Fmt.Sprintf("?language=%s", language)
	}
	return webIntegration.Get(webIntegration.forURL(baseUrl))
}

func (webIntegration *WebIntegration) ManagementPostTranslations(createCms *Management.CreateCms) *WebIntegration {
	return webIntegration.Post(webIntegration.forURL("/api/management/v1/translations"), createCms)
}

func (webIntegration *WebIntegration) ProductManagementGetProducts() *WebIntegration {
	return webIntegration.Get(webIntegration.forURL("/api/product-management/v1/products"))
}

func (webIntegration *WebIntegration) ShoppingBasketCreate() *WebIntegration {
	return webIntegration.Post(webIntegration.forURL("/api/shopping-basket/v1/shopping-baskets"), nil)
}

func (webIntegration *WebIntegration) ShoppingBasketAddItem(ID string, updateShoppingBasketItem ShoppingBasket.UpdateShoppingBasketItem) *WebIntegration {
	return webIntegration.Post(webIntegration.forURL("/api/shopping-basket/v1/shopping-baskets/"+ID), updateShoppingBasketItem)
}

func (webIntegration *WebIntegration) GetShoppingBasket(ID string) *WebIntegration {
	return webIntegration.Get(webIntegration.forURL("/api/shopping-basket/v1/shopping-baskets/" + ID))
}

func (webIntegration *WebIntegration) ProductManagementGetProductByID(ID string) *WebIntegration {
	return webIntegration.Get(webIntegration.forURL("/api/product-management/v1/products/" + ID))
}

func (webIntegration *WebIntegration) forURL(URL string) string {
	return webIntegration.server.URL + URL
}

func (webIntegration *WebIntegration) DatabaseConnection() *Database.DB {
	return webIntegration.databaseConnection
}

func (webIntegration *WebIntegration) Get(URL string) *WebIntegration {
	return webIntegration.doRequest("GET", URL, nil)
}

func (webIntegration *WebIntegration) Delete(URL string) *WebIntegration {
	return webIntegration.doRequest("DELETE", URL, nil)
}

func (webIntegration *WebIntegration) Post(URL string, body any) *WebIntegration {
	return webIntegration.doRequest("POST", URL, body)
}

func (webIntegration *WebIntegration) Put(URL string, body any) *WebIntegration {
	return webIntegration.doRequest("PUT", URL, body)
}

func (webIntegration *WebIntegration) GetResponseBody(body any) *WebIntegration {
	err := JSON.NewDecoder(webIntegration.response.Body).Decode(&body)
	Assert.Nil(webIntegration.unitTest, err)
	return webIntegration
}

func (webIntegration *WebIntegration) AssertStatusCreated() *WebIntegration {
	return webIntegration.Equal(Http.StatusCreated, webIntegration.response.StatusCode)
}

func (webIntegration *WebIntegration) IsNotNil(body any) *WebIntegration {
	Assert.NotNil(webIntegration.unitTest, body)
	return webIntegration
}

func (webIntegration *WebIntegration) AssertStatusOk() *WebIntegration {
	return webIntegration.Equal(Http.StatusOK, webIntegration.response.StatusCode)
}

func (webIntegration *WebIntegration) Equal(expected, actual any, arguments ...any) *WebIntegration {
	Assert.Equal(webIntegration.unitTest, expected, actual, arguments...)
	return webIntegration
}

func (webIntegration *WebIntegration) AssertBadRequest() *WebIntegration {
	return webIntegration.Equal(Http.StatusBadRequest, webIntegration.response.StatusCode)
}

func (webIntegration *WebIntegration) doRequest(method string, URL string, body any) *WebIntegration {
	var reader io.Reader
	if body != nil {
		serializedBody, err := JSON.Marshal(body)
		Assert.Nil(webIntegration.unitTest, err)
		reader = bytes.NewBuffer(serializedBody)
	}
	request, err := Http.NewRequest(method, URL, reader)
	Assert.Nil(webIntegration.unitTest, err)
	if err != nil || request == nil {
		Assert.Fail(webIntegration.unitTest, "Could not create the request")
		return webIntegration
	}
	request.Header.Set("Content-Type", "application/json")

	if webIntegration.authToken != "" {
		request.Header.Set(ApplicationMiddleware.Authorization, Fmt.Sprintf("Bearer %s", webIntegration.authToken))
	}

	client := &Http.Client{}

	webIntegration.response, err = client.Do(request)

	Assert.Nil(webIntegration.unitTest, err)

	return webIntegration
}

func createTestUser() *ApplicationMiddleware.UserClaims {
	return &ApplicationMiddleware.UserClaims{
		UserID: "UserId",
		Roles:  []string{"ADMIN", "USER"},
	}
}

func SetupWebIntegration(unitTest *Testing.T) *WebIntegration {
	var token string
	secretKeyProvider := ApplicationMiddleware.NewSecretKeyProvider()
	secretKey, err := secretKeyProvider.GenerateKey()
	if err != nil {
		Log.Println(err.Error())
	}

	unitTest.Helper()
	backendInitializer := NewBackendInitializer()
	backendInitializer.TestContainers(unitTest)
	backendInitializer.Run()
	databaseConnection := backendInitializer.DatabaseConnection()
	serverConfiguration := Configuration.NewServerConfiguration()
	serverConfiguration.JwtSecret = secretKey
	databaseConnectionFactory := Factory.NewDatabaseConnectionFactory(backendInitializer.DatabaseConfiguration)
	messageBrokerConnectionFactory := Factory.NewMessageBrokerConnectionFactory(backendInitializer.MessageBrokerConfiguration)
	applicationConnectionFactory := Factory.NewApplicationConnectionFactory(databaseConnectionFactory, messageBrokerConnectionFactory)

	ShoppingBasket.RegisterShoppingBasketEvents(applicationConnectionFactory.ShoppingBasketService(), messageBrokerConnectionFactory.MessageBroker())
	if err := messageBrokerConnectionFactory.MessageBroker().Start(); err != nil {
		Log.Println(err.Error())
	}
	apiDefinition := ApplicationRouter.APIDefinition{
		ServerConfiguration:          &serverConfiguration,
		ApplicationConnectionFactory: applicationConnectionFactory,
	}

	server := HttpTest.NewServer(apiDefinition.ConfigRouter())

	unitTest.Cleanup(func() {
		server.Close()
	})
	token, err = ApplicationMiddleware.EncryptClaims(createTestUser(), secretKey)
	if err != nil {
		unitTest.Fatal(err)
	}

	return &WebIntegration{
		databaseConnection: databaseConnection,
		server:             server,
		unitTest:           unitTest,
		authToken:          token,
	}
}
