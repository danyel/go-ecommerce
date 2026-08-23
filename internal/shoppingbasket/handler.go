package shoppingbasket

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

type ShoppingBasketHandler interface {
	CreateShoppingBasket(response Http.ResponseWriter, request *Http.Request)
	UpdateShoppingBasketItem(response Http.ResponseWriter, request *Http.Request)
	GetShoppingBasket(response Http.ResponseWriter, request *Http.Request)
}

type shoppingBasketHandler struct {
	shoppingBasketService ShoppingBasketService
}

func (shoppingBasketHandler *shoppingBasketHandler) CreateShoppingBasket(response Http.ResponseWriter, request *Http.Request) {
	shoppingBasket, err := shoppingBasketHandler.shoppingBasketService.CreateShoppingBasket()
	if err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, shoppingBasket.ID)
}

// UpdateShoppingBasketItem web handler function that will update the shopping basket
func (shoppingBasketHandler *shoppingBasketHandler) UpdateShoppingBasketItem(response Http.ResponseWriter, request *Http.Request) {
	var updateShoppingBasketItem UpdateShoppingBasketItem
	var err error
	var shoppingBasket ShoppingBasket
	shoppingBasketID, err := WebHandler.GetID(request, "shoppingBasketID")
	if err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	if err = WebHandler.ValidateRequest(request, &updateShoppingBasketItem); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if shoppingBasket, err = shoppingBasketHandler.shoppingBasketService.UpdateShoppingBasketItem(shoppingBasketID.ID, updateShoppingBasketItem); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func (shoppingBasketHandler *shoppingBasketHandler) GetShoppingBasket(response Http.ResponseWriter, request *Http.Request) {
	var err error
	var shoppingBasket ShoppingBasket
	ID, err := WebHandler.GetID(request, "shoppingBasketID")
	if err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if shoppingBasket, err = shoppingBasketHandler.shoppingBasketService.GetShoppingBasket(ID.ID); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func NewHandler(shoppingBasketService ShoppingBasketService) ShoppingBasketHandler {
	return &shoppingBasketHandler{shoppingBasketService}
}
