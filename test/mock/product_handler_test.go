package mock

import (
	Http "net/http"
	Testing "testing"

	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	Uuid "github.com/google/uuid"
	Assert "github.com/stretchr/testify/assert"
	Mock "github.com/stretchr/testify/mock"
)

type MockProductService struct {
	Mock.Mock
}

func (m *MockProductService) GetProducts() []Product.Product {
	args := m.Called()
	return args.Get(0).([]Product.Product)
}

func (m *MockProductService) GetProduct(uuid Uuid.UUID) (Product.Product, error) {
	args := m.Called(uuid)
	return args.Get(0).(Product.Product), args.Error(1)
}

func TestProductHandler(t *Testing.T) {
	mockProductService := new(MockProductService)
	h := Product.NewAPIHandler(mockProductService)
	mh := Run(t)

	t.Run("GetProducts", func(t *Testing.T) {
		products := []Product.Product{
			{
				Code:  "Code",
				Price: Types.NewPrice(1000, "EUR"),
			},
		}
		mockProductService.On("GetProducts").Return(products, nil)
		Assert.Equal(t, Http.StatusOK, mh.New().
			NewRecoder().
			NewRequest(Http.MethodGet, "/api/product/v1/products", nil).
			NewRouter(Http.MethodGet, "/api/product/v1/products", h.GetProducts).
			ServeHTTP().
			Status())
		mockProductService.AssertCalled(t, "GetProducts")
		mockProductService.AssertExpectations(t)
	})

	t.Run("GetProduct", func(t *Testing.T) {
		i, _ := Uuid.Parse("aef8f0ce-c33f-456c-bc5c-91f951116cf7")
		p := Product.Product{Code: "Code", Price: Types.NewPrice(1000, "EUR")}
		mockProductService.On("GetProduct", i).Return(p, nil)

		Assert.Equal(t, Http.StatusOK, mh.New().
			NewRecoder().
			NewRequest(Http.MethodGet, "/api/product/v1/products/"+i.String(), nil).
			NewRouter(Http.MethodGet, "/api/product/v1/products/{productID}", h.GetProduct).
			ServeHTTP().
			Status())
		mockProductService.AssertCalled(t, "GetProduct", i)
		mockProductService.AssertExpectations(t)
	})
}
