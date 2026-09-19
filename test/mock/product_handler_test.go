package mock

import (
	Http "net/http"
	Testing "testing"

	ApplicationRouter "github.com/danyel/ecommerce/cmd/router"
	Domain "github.com/danyel/ecommerce/internal/domain"
	Handler "github.com/danyel/ecommerce/internal/handler"
	Mapper "github.com/danyel/ecommerce/internal/mapper"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Types "github.com/danyel/ecommerce/internal/types"
	SetupWebIntegration "github.com/danyel/ecommerce/test/integration/initializer"
	TestUtils "github.com/danyel/ecommerce/test/testutils"
	Uuid "github.com/google/uuid"
	Assert "github.com/stretchr/testify/assert"
	Mock "github.com/stretchr/testify/mock"
)

type MockProductService struct {
	Mock.Mock
}

type MockProductMapper struct {
	Mock.Mock
}

func (productService *MockProductService) FindAll() []Domain.Product {
	args := productService.Called()
	return args.Get(0).([]Domain.Product)
}

func (productService *MockProductService) FindByID(ID Uuid.UUID) (Domain.Product, error) {
	args := productService.Called(ID)
	return args.Get(0).(Domain.Product), args.Error(1)
}

func (productService *MockProductService) Update(product Domain.Product) error {
	args := productService.Called(product)
	return args.Get(0).(error)
}

func (productService *MockProductService) UpdateStock(ID Uuid.UUID, shoppingBasketId Uuid.UUID, stock int) error {
	args := productService.Called(ID, shoppingBasketId, stock)
	return args.Get(0).(error)
}

func (productMapper *MockProductMapper) MapProducts(models []*Persistence.ProductModel) []Domain.Product {
	args := productMapper.Called(models)
	return args.Get(0).([]Domain.Product)
}

func (productMapper *MockProductMapper) MapProduct(model *Persistence.ProductModel) Domain.Product {
	args := productMapper.Called(model)
	return args.Get(0).(Domain.Product)
}

func TestProductHandler(unitTest *Testing.T) {
	TestUtils.PreInitTest()
	productService := new(MockProductService)
	productMapper := new(MockProductMapper)
	productHandler := Handler.ProductWebHandler(productService, productMapper, Mapper.CategoryMapper())
	run := Run(unitTest)

	unitTest.Run("FindAll", func(unitTest *Testing.T) {
		products := []Domain.Product{
			{
				Code:  "Code",
				Price: Types.NewPrice(1000, "EUR"),
			},
		}
		productService.On("FindAll").Return(products, nil)
		Assert.Equal(unitTest, Http.StatusOK, run.New().
			NewRecoder().
			NewRequest(Http.MethodGet, SetupWebIntegration.ProductProductsURL, nil).
			NewRouter(Http.MethodGet, SetupWebIntegration.ProductProductsURL, productHandler.HandleGetProductsV1).
			ServeHTTP().
			Status())
		productService.AssertCalled(unitTest, "FindAll")
		productService.AssertExpectations(unitTest)
	})

	unitTest.Run("FindByID", func(unitTest *Testing.T) {
		ID, _ := Uuid.Parse("aef8f0ce-c33f-456c-bc5c-91f951116cf7")
		product := Domain.Product{Code: "Code", Price: Types.NewPrice(1000, "EUR")}
		productService.On("FindByID", ID).Return(product, nil)

		Assert.Equal(unitTest, Http.StatusOK, run.New().
			NewRecoder().
			NewRequest(Http.MethodGet, SetupWebIntegration.ProductProductsURL+ApplicationRouter.SLASH+ID.String(), nil).
			NewRouter(Http.MethodGet, SetupWebIntegration.ProductProductsURL+ApplicationRouter.ByID, productHandler.HandleGetProductV1).
			ServeHTTP().
			Status())
		productService.AssertCalled(unitTest, "FindByID", ID)
		productService.AssertExpectations(unitTest)
	})
}
