package shoppingbasket

import (
	Fmt "fmt"
	"net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	Uuid "github.com/google/uuid"
)

type Validator interface {
	Validate(ID Types.ID) ShoppingBasket
	ValidateItem(ID Types.ID, updateShoppingBasketItem UpdateShoppingBasketItemDTO) WebHandler.ProblemDetail
}

type validator struct {
	productRepository        Repository.CrudRepository[Product.ProductModel]
	shoppingBasketRepository Repository.CrudRepository[ShoppingBasketModel]
	productService           Product.ProductService
}

func (validator *validator) Validate(ID Types.ID) ShoppingBasket {
	var shoppingBasket = ShoppingBasket{
		ProblemDetail: WebHandler.ProblemDetail{
			Status: 0,
		},
	}
	if ID.ID == Uuid.Nil {
		shoppingBasket.ProblemDetail.Status = http.StatusBadRequest
		shoppingBasket.ProblemDetail.Errors = make(map[string]any)
		shoppingBasket.ProblemDetail.Errors["ID"] = Fmt.Sprintf("ID is invalid: %v", ID.ID)
		return shoppingBasket
	}
	var shoppingBasketModel *ShoppingBasketModel
	shoppingBasketModel, err := validator.shoppingBasketRepository.FindByID(ID.ID, "Items")
	if err != nil {
		details := make(map[string]any)
		details["id"] = Fmt.Sprintf("Could not find shopping basket with id '%s'", ID.ID.String())
		shoppingBasket.ProblemDetail.Status = http.StatusNotFound
		shoppingBasket.ProblemDetail.Title = WebHandler.NotFoundTitle
		shoppingBasket.ProblemDetail.Errors = details
		return shoppingBasket
	}
	shoppingBasket.ID = Types.NewID(shoppingBasketModel.ID)
	shoppingBasket.Items = make([]ShoppingBasketItem, len(shoppingBasketModel.Items))
	for i, item := range shoppingBasketModel.Items {
		product, err := validator.productService.GetProduct(item.ProductID)
		if err != nil {
			details := make(map[string]any, 1)
			details["product_id"] = Fmt.Sprintf("Product not found: '%s'", item.ProductID.String())
			shoppingBasket.ProblemDetail.Status = http.StatusNotFound
			shoppingBasket.ProblemDetail.Title = WebHandler.NotFoundTitle
			shoppingBasket.ProblemDetail.Errors = details
		} else {
			shoppingBasket.Items[i] = ShoppingBasketItem{
				product,
				item.Quantity,
			}
		}
	}
	return shoppingBasket
}

func (validator *validator) ValidateItem(ID Types.ID, updateShoppingBasketItem UpdateShoppingBasketItemDTO) WebHandler.ProblemDetail {
	if _, err := validator.shoppingBasketRepository.FindByID(ID.ID, "Items"); err != nil {
		details := make(map[string]any, 1)
		details["id"] = Fmt.Sprintf("Shopping basket not found: '%s'", ID.ID.String())
		return WebHandler.ProblemDetail{
			Title:  WebHandler.NotFoundTitle,
			Status: http.StatusNotFound,
			Errors: details,
		}
	}
	product, err := validator.productRepository.FindByID(updateShoppingBasketItem.ProductID.ID)
	if err != nil {
		details := make(map[string]any, 1)
		details["product_id"] = Fmt.Sprintf("Product not found: '%s'", updateShoppingBasketItem.ProductID.ID.String())
		return WebHandler.ProblemDetail{
			Title:  WebHandler.NotFoundTitle,
			Status: http.StatusNotFound,
			Errors: details,
		}
	}

	if product.Stock == 0 && updateShoppingBasketItem.Quantity > 0 {
		details := make(map[string]any, 1)
		details["product_stock"] = Fmt.Sprintf("'%s' is out of stock", product.Name)
		return WebHandler.ProblemDetail{
			Title:  WebHandler.BadRequestTitle,
			Status: http.StatusBadRequest,
			Errors: details,
		}

	}

	if product.Stock != 0 && product.Stock < updateShoppingBasketItem.Quantity {
		details := make(map[string]any, 1)
		details["product_stock"] = Fmt.Sprintf("'%s' can not reserve all items: %d/%d", product.Name, product.Stock, updateShoppingBasketItem.Quantity)
		return WebHandler.ProblemDetail{
			Title:  WebHandler.BadRequestTitle,
			Status: http.StatusBadRequest,
			Errors: details,
		}
	}

	return WebHandler.ProblemDetail{Status: 0}
}

func NewValidator(productRepository Repository.CrudRepository[Product.ProductModel], shoppingBasketRepository Repository.CrudRepository[ShoppingBasketModel], productService Product.ProductService) Validator {
	return &validator{productRepository, shoppingBasketRepository, productService}
}
