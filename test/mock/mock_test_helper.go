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
	w *HttpTest.ResponseRecorder
	r *Http.Request
	m *Router.Mux
}

func (h *MockHelper) New() MockFluent {
	return MockFluent{}
}

func (m MockFluent) NewRecoder() MockFluent {
	m.w = HttpTest.NewRecorder()
	return m
}

func (m MockFluent) NewRouter(method string, pattern string, handlerFn Http.HandlerFunc) MockFluent {
	m.m = Router.NewRouter()
	m.m.Method(method, pattern, handlerFn)
	return m
}

func (m MockFluent) NewRequest(method string, target string, body IO.Reader) MockFluent {
	m.r = HttpTest.NewRequest(method, target, body)
	return m
}

func (m MockFluent) ServeHTTP() MockFluent {
	m.m.ServeHTTP(m.w, m.r)
	return m
}

func (m MockFluent) Status() int {
	return m.w.Result().StatusCode
}

func Run(t *Testing.T) *MockHelper {
	t.Helper()
	return &MockHelper{}
}
