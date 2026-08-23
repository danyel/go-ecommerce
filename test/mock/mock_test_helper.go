package mock

import (
	IO "io"
	Http "net/http"
	HttpTest "net/http/httptest"
	Testing "testing"

	Router "github.com/go-chi/chi/v5"
	Mock "github.com/stretchr/testify/mock"
)

//goland:noinspection GoNameStartsWithPackageName
type MockHelper struct {
	Mock.Mock
}

//goland:noinspection GoNameStartsWithPackageName
type MockFluent struct {
	response *HttpTest.ResponseRecorder
	request  *Http.Request
	router   *Router.Mux
}

func (mockHelper *MockHelper) New() MockFluent {
	return MockFluent{}
}

func (mockFluent MockFluent) NewRecoder() MockFluent {
	mockFluent.response = HttpTest.NewRecorder()
	return mockFluent
}

func (mockFluent MockFluent) NewRouter(method string, pattern string, handlerFunc Http.HandlerFunc) MockFluent {
	mockFluent.router = Router.NewRouter()
	mockFluent.router.Method(method, pattern, handlerFunc)
	return mockFluent
}

func (mockFluent MockFluent) NewRequest(method string, target string, reader IO.Reader) MockFluent {
	mockFluent.request = HttpTest.NewRequest(method, target, reader)
	return mockFluent
}

func (mockFluent MockFluent) ServeHTTP() MockFluent {
	mockFluent.router.ServeHTTP(mockFluent.response, mockFluent.request)
	return mockFluent
}

func (mockFluent MockFluent) Status() int {
	return mockFluent.response.Result().StatusCode
}

func Run(unitTest *Testing.T) *MockHelper {
	unitTest.Helper()
	return &MockHelper{}
}
