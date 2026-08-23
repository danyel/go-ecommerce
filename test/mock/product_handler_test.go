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

func (productService *MockProductService) GetProducts() []Product.Product {
	args := productService.Called()
	return args.Get(0).([]Product.Product)
}

func (productService *MockProductService) GetProduct(ID Uuid.UUID) (Product.Product, error) {
	args := productService.Called(ID)
	return args.Get(0).(Product.Product), args.Error(1)
}

func TestProductHandler(unitTest *Testing.T) {
	productService := new(MockProductService)
	productHandler := Product.NewHandler(productService)
	run := Run(unitTest)

	unitTest.Run("GetProducts", func(unitTest *Testing.T) {
		products := []Product.Product{
			{
				Code:  "Code",
				Price: Types.NewPrice(1000, "EUR"),
			},
		}
		productService.On("GetProducts").Return(products, nil)
		Assert.Equal(unitTest, Http.StatusOK, run.New().
			NewRecoder().
			NewRequest(Http.MethodGet, "/api/product/v1/products", nil).
			NewRouter(Http.MethodGet, "/api/product/v1/products", productHandler.GetProducts).
			ServeHTTP().
			Status())
		productService.AssertCalled(unitTest, "GetProducts")
		productService.AssertExpectations(unitTest)
	})

	unitTest.Run("GetProduct", func(unitTest *Testing.T) {
		ID, _ := Uuid.Parse("aef8f0ce-c33f-456c-bc5c-91f951116cf7")
		product := Product.Product{Code: "Code", Price: Types.NewPrice(1000, "EUR")}
		productService.On("GetProduct", ID).Return(product, nil)

		Assert.Equal(unitTest, Http.StatusOK, run.New().
			NewRecoder().
			NewRequest(Http.MethodGet, "/api/product/v1/products/"+ID.String(), nil).
			NewRouter(Http.MethodGet, "/api/product/v1/products/{productID}", productHandler.GetProduct).
			ServeHTTP().
			Status())
		productService.AssertCalled(unitTest, "GetProduct", ID)
		productService.AssertExpectations(unitTest)
	})
}
