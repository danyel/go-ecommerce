package shoppingbasket

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketWebHandler interface {
	HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request)
	HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetShoppingBasketV1(response Http.ResponseWriter, request *Http.Request)
}

type shoppingBasketWebHandler struct {
	shoppingBasketService ShoppingBasketService
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request) {
	shoppingBasket, err := shoppingBasketWebHandler.shoppingBasketService.CreateShoppingBasket()
	if err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, shoppingBasket.ID)
}

// HandleUpdateShoppingBasketItemV1 web handler function that will update the shopping basket
func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request) {
	var updateShoppingBasketItem UpdateShoppingBasketItem
	var err error
	var shoppingBasket ShoppingBasket
	shoppingBasketID, err := WebHandler.GetID(request)
	if err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	if err = WebHandler.ValidateRequest(request, &updateShoppingBasketItem); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.UpdateShoppingBasketItem(shoppingBasketID.ID, updateShoppingBasketItem); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleGetShoppingBasketV1(response Http.ResponseWriter, request *Http.Request) {
	var err error
	var shoppingBasket ShoppingBasket
	ID, err := WebHandler.GetID(request)
	if err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.GetShoppingBasket(ID.ID); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func NewWebHandler(shoppingBasketService ShoppingBasketService) ShoppingBasketWebHandler {
	return &shoppingBasketWebHandler{shoppingBasketService}
}
