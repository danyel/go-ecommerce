package shoppingbasket

import (
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketWebHandler interface {
	HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request)
	HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetShoppingBasketByIdV1(response Http.ResponseWriter, request *Http.Request)
}

type shoppingBasketWebHandler struct {
	shoppingBasketService ShoppingBasketService
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleCreateShoppingBasketV1")
	shoppingBasket, err := shoppingBasketWebHandler.shoppingBasketService.CreateShoppingBasket()
	if err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, shoppingBasket.ID)
}

// HandleUpdateShoppingBasketItemV1 web handler function that will update the shopping basket
func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleUpdateShoppingBasketItemV1")
	var updateShoppingBasketItem UpdateShoppingBasketItem
	var err error
	var shoppingBasket ShoppingBasket
	ID, err := WebHandler.GetID(request)
	if err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}

	if err = WebHandler.ValidateRequest(request, &updateShoppingBasketItem); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.UpdateShoppingBasketItem(ID.ID, updateShoppingBasketItem); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleGetShoppingBasketByIdV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleGetShoppingBasketByIdV1")
	var err error
	var shoppingBasket ShoppingBasket
	ID, err := WebHandler.GetID(request)
	if err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	Logger.Log.Debug("Fetching for id: %s", ID.ID.String())
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.GetShoppingBasket(ID.ID); err != nil {
		Logger.Log.Debug(err.Error())
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	Logger.Log.Debug("Shopping Basket fetched: %v", shoppingBasket)
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func NewWebHandler(shoppingBasketService ShoppingBasketService) ShoppingBasketWebHandler {
	return &shoppingBasketWebHandler{shoppingBasketService}
}
