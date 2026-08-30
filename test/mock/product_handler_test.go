package mock

import (
	Http "net/http"
	Testing "testing"

	ApplicationRouter "github.com/danyel/ecommerce/cmd/router"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	SetupWebIntegration "github.com/danyel/ecommerce/test/integration/initializer"
	TestUtils "github.com/danyel/ecommerce/test/testutils"
	Uuid "github.com/google/uuid"
	Assert "github.com/stretchr/testify/assert"
	Mock "github.com/stretchr/testify/mock"
)

type MockProductService struct {
	Mock.Mock
}

func (productService *MockProductService) GetProducts() []Product.ProductDTO {
	args := productService.Called()
	return args.Get(0).([]Product.ProductDTO)
}

func (productService *MockProductService) GetProduct(ID Uuid.UUID) (Product.ProductDTO, error) {
	args := productService.Called(ID)
	return args.Get(0).(Product.ProductDTO), args.Error(1)
}

func TestProductHandler(unitTest *Testing.T) {
	TestUtils.PreInitTest()
	productService := new(MockProductService)
	productHandler := Product.NewWebHandler(productService)
	run := Run(unitTest)

	unitTest.Run("GetProducts", func(unitTest *Testing.T) {
		products := []Product.ProductDTO{
			{
				Code:  "Code",
				Price: Types.NewPrice(1000, "EUR"),
			},
		}
		productService.On("GetProducts").Return(products, nil)
		Assert.Equal(unitTest, Http.StatusOK, run.New().
			NewRecoder().
			NewRequest(Http.MethodGet, SetupWebIntegration.ProductProductsURL, nil).
			NewRouter(Http.MethodGet, SetupWebIntegration.ProductProductsURL, productHandler.HandleGetProductsV1).
			ServeHTTP().
			Status())
		productService.AssertCalled(unitTest, "GetProducts")
		productService.AssertExpectations(unitTest)
	})

	unitTest.Run("GetProduct", func(unitTest *Testing.T) {
		ID, _ := Uuid.Parse("aef8f0ce-c33f-456c-bc5c-91f951116cf7")
		product := Product.ProductDTO{Code: "Code", Price: Types.NewPrice(1000, "EUR")}
		productService.On("GetProduct", ID).Return(product, nil)

		Assert.Equal(unitTest, Http.StatusOK, run.New().
			NewRecoder().
			NewRequest(Http.MethodGet, SetupWebIntegration.ProductProductsURL+ApplicationRouter.SLASH+ID.String(), nil).
			NewRouter(Http.MethodGet, SetupWebIntegration.ProductProductsURL+ApplicationRouter.ByID, productHandler.HandleGetProductV1).
			ServeHTTP().
			Status())
		productService.AssertCalled(unitTest, "GetProduct", ID)
		productService.AssertExpectations(unitTest)
	})
}
