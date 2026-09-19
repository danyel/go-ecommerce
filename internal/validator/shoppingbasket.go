package validator

import (
	Fmt "fmt"
	Http "net/http"

	Domain "github.com/danyel/ecommerce/internal/domain"
	Model "github.com/danyel/ecommerce/internal/model"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Types "github.com/danyel/ecommerce/internal/types"
	Uuid "github.com/google/uuid"
)

type Validator interface {
	Validate(ID Types.ID) Domain.ShoppingBasket
	ValidateItem(ID Types.ID, updateShoppingBasketItem Model.UpdateShoppingBasketItemDTO) Domain.ProblemDetail
}

type productFinder interface {
	FindByID(Uuid.UUID) (Domain.Product, error)
}

type validator struct {
	productRepository        Persistence.CrudRepository[Persistence.ProductModel]
	shoppingBasketRepository Persistence.CrudRepository[Persistence.ShoppingBasketModel]
	productService           productFinder
}

func (validator *validator) Validate(ID Types.ID) Domain.ShoppingBasket {
	var shoppingBasket = Domain.ShoppingBasket{
		ProblemDetail: Domain.ProblemDetail{
			Status: 0,
		},
	}
	if ID.ID == Uuid.Nil {
		shoppingBasket.ProblemDetail.Status = Http.StatusBadRequest
		shoppingBasket.ProblemDetail.Errors = make(map[string]any)
		shoppingBasket.ProblemDetail.Errors["ID"] = Fmt.Sprintf("ID is invalid: %v", ID.ID)
		return shoppingBasket
	}
	var shoppingBasketModel *Persistence.ShoppingBasketModel
	shoppingBasketModel, err := validator.shoppingBasketRepository.FindByID(ID.ID, "Items")
	if err != nil {
		details := make(map[string]any)
		details["id"] = Fmt.Sprintf("Could not find shopping basket with id '%s'", ID.ID.String())
		shoppingBasket.ProblemDetail.Status = Http.StatusNotFound
		shoppingBasket.ProblemDetail.Title = "Not Found"
		shoppingBasket.ProblemDetail.Errors = details
		return shoppingBasket
	}
	shoppingBasket.ID = Types.NewID(shoppingBasketModel.ID)
	shoppingBasket.Items = make([]Domain.ShoppingBasketItem, len(shoppingBasketModel.Items))
	for i, item := range shoppingBasketModel.Items {
		product, err := validator.productService.FindByID(item.ProductID)
		if err != nil {
			details := make(map[string]any, 1)
			details["product_id"] = Fmt.Sprintf("Product not found: '%s'", item.ProductID.String())
			shoppingBasket.ProblemDetail.Status = Http.StatusNotFound
			shoppingBasket.ProblemDetail.Title = "Not Found"
			shoppingBasket.ProblemDetail.Errors = details
		} else {
			shoppingBasket.Items[i] = Domain.ShoppingBasketItem{
				Product:  product,
				Quantity: item.Quantity,
			}
		}
	}
	return shoppingBasket
}

func (validator *validator) ValidateItem(ID Types.ID, updateShoppingBasketItem Model.UpdateShoppingBasketItemDTO) Domain.ProblemDetail {
	if updateShoppingBasketItem.Quantity < 0 {
		return Domain.ProblemDetail{
			Title:  "Bad Request",
			Status: Http.StatusBadRequest,
			Errors: map[string]any{
				"quantity": "quantity cannot be negative",
			},
		}
	}

	shoppingBasketModel, err := validator.shoppingBasketRepository.FindByID(ID.ID, "Items")
	if err != nil {
		details := make(map[string]any, 1)
		details["id"] = Fmt.Sprintf("Shopping basket not found: '%s'", ID.ID.String())
		return Domain.ProblemDetail{
			Title:  "Not Found",
			Status: Http.StatusNotFound,
			Errors: details,
		}
	}
	product, err := validator.productRepository.FindByID(updateShoppingBasketItem.ProductID.ID)
	if err != nil {
		details := make(map[string]any, 1)
		details["product_id"] = Fmt.Sprintf("Product not found: '%s'", updateShoppingBasketItem.ProductID.ID.String())
		return Domain.ProblemDetail{
			Title:  "Not Found",
			Status: Http.StatusNotFound,
			Errors: details,
		}
	}

	currentQuantity := 0
	for _, item := range shoppingBasketModel.Items {
		if item.ProductID == product.ID {
			currentQuantity = item.Quantity
			break
		}
	}
	availableStock := product.Stock + currentQuantity
	if availableStock == 0 && updateShoppingBasketItem.Quantity > 0 {
		details := make(map[string]any, 1)
		details["product_stock"] = Fmt.Sprintf("'%s' is out of stock", product.Name)
		return Domain.ProblemDetail{
			Title:  "Not Found",
			Status: Http.StatusBadRequest,
			Errors: details,
		}

	}

	if updateShoppingBasketItem.Quantity > availableStock {
		details := make(map[string]any, 1)
		details["product_stock"] = Fmt.Sprintf("'%s' can not reserve all items: %d/%d", product.Name, product.Stock, updateShoppingBasketItem.Quantity)
		return Domain.ProblemDetail{
			Title:  "Bad Request",
			Status: Http.StatusBadRequest,
			Errors: details,
		}
	}

	return Domain.ProblemDetail{Status: 0}
}

func NewValidator(productRepository Persistence.CrudRepository[Persistence.ProductModel], shoppingBasketRepository Persistence.CrudRepository[Persistence.ShoppingBasketModel], productService productFinder) Validator {
	return &validator{productRepository, shoppingBasketRepository, productService}
}
