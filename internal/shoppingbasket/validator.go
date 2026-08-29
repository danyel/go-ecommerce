package shoppingbasket

import (
	Errors "errors"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Service "github.com/danyel/ecommerce/internal/common/service"
	Product "github.com/danyel/ecommerce/internal/product"
)

const (
	OutOfStock        Service.Status = "OUT_OF_STOCK"
	NotFound          Service.Status = "NOT_FOUND"
	OutOfStockMessage string         = "Product Out of Stock"
	NotFoundMessage   string         = "Product Not Found"
)

var ErrorOutOfStock = Errors.New("no more items in stock")

type Validator interface {
	ValidateItem(updateShoppingBasketItem UpdateShoppingBasketItem) error
}

type validator struct {
	productRepository Repository.CrudRepository[Product.ProductModel]
}

func (validator *validator) ValidateItem(updateShoppingBasketItem UpdateShoppingBasketItem) error {
	Logger.Log.Debug("validate shopping basket item: %v", updateShoppingBasketItem)
	product, err := validator.productRepository.FindById(updateShoppingBasketItem.ProductID.ID)
	if err != nil {
		errs := make(map[string]string, 1)
		errs["id"] = updateShoppingBasketItem.ProductID.ID.String()
		return &Service.ServiceError{
			Status:  NotFound,
			Message: NotFoundMessage,
			Fields:  errs,
		}
	}

	if product.Stock == 0 && updateShoppingBasketItem.Quantity > 0 {
		errs := make(map[string]string, 1)
		errs["id"] = updateShoppingBasketItem.ProductID.ID.String()
		return &Service.ServiceError{
			Status:  OutOfStock,
			Message: OutOfStockMessage,
			Fields:  errs,
		}
	}

	return nil
}

func NewValidator(productRepository Repository.CrudRepository[Product.ProductModel]) Validator {
	return &validator{productRepository}
}
