package mock

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

type MockHelper struct {
	mock.Mock
}

type MockFluent struct {
	w *httptest.ResponseRecorder
	r *http.Request
	m *chi.Mux
}

func (h *MockHelper) New() MockFluent {
	return MockFluent{}
}

func (m MockFluent) NewRecoder() MockFluent {
	m.w = httptest.NewRecorder()
	return m
}

func (m MockFluent) NewRouter(method string, pattern string, handlerFn http.HandlerFunc) MockFluent {
	m.m = chi.NewRouter()
	m.m.Method(method, pattern, handlerFn)
	return m
}

func (m MockFluent) NewRequest(method string, target string, body io.Reader) MockFluent {
	m.r = httptest.NewRequest(method, target, body)
	return m
}

func (m MockFluent) ServeHTTP() MockFluent {
	m.m.ServeHTTP(m.w, m.r)
	return m
}

func (m MockFluent) Status() int {
	return m.w.Result().StatusCode
}

func Run(t *testing.T) *MockHelper {
	t.Helper()
	return &MockHelper{}
}
