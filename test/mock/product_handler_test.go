package mock

import (
	"net/http"
	"testing"

	"github.com/dnoulet/ecommerce/internal/product"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProducts() []product.Product {
	args := m.Called()
	return args.Get(0).([]product.Product)
}

func (m *MockProductService) GetProduct(uuid uuid.UUID) (product.Product, error) {
	args := m.Called(uuid)
	return args.Get(0).(product.Product), args.Error(1)
}

func TestProductHandler(t *testing.T) {
	mockProductService := new(MockProductService)
	h := product.NewHandler(mockProductService)
	mh := Run(t)

	t.Run("GetProducts", func(t *testing.T) {
		products := []product.Product{
			{
				Code:  "Code",
				Price: 1000,
			},
		}
		mockProductService.On("GetProducts").Return(products, nil)
		assert.Equal(t, http.StatusOK, mh.New().
			NewRecoder().
			NewRequest(http.MethodGet, "/api/product/v1/products", nil).
			NewRouter(http.MethodGet, "/api/product/v1/products", h.GetProducts).
			ServeHTTP().
			Status())
		mockProductService.AssertCalled(t, "GetProducts")
		mockProductService.AssertExpectations(t)
	})

	t.Run("GetProduct", func(t *testing.T) {
		i, _ := uuid.Parse("aef8f0ce-c33f-456c-bc5c-91f951116cf7")
		p := product.Product{Code: "Code", Price: 1000}
		mockProductService.On("GetProduct", i).Return(p, nil)

		assert.Equal(t, http.StatusOK, mh.New().
			NewRecoder().
			NewRequest(http.MethodGet, "/api/product/v1/products/"+i.String(), nil).
			NewRouter(http.MethodGet, "/api/product/v1/products/{productId}", h.GetProduct).
			ServeHTTP().
			Status())
		mockProductService.AssertCalled(t, "GetProduct", i)
		mockProductService.AssertExpectations(t)
	})
}
