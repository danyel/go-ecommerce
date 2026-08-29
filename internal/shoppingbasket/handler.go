package shoppingbasket

import (
	"errors"
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
)

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketWebHandler interface {
	HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request)
	HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetShoppingBasketByIDV1(response Http.ResponseWriter, request *Http.Request)
}

type shoppingBasketWebHandler struct {
	shoppingBasketService ShoppingBasketService
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleCreateShoppingBasketV1")
	shoppingBasket, err := shoppingBasketWebHandler.shoppingBasketService.CreateShoppingBasket()
	if err != nil {
		WebHandler.InternalServerError(response, request, "Error while creating shopping basket", err.Error())
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
		WebHandler.BadRequest(response, request, "ID not found", err.Error())
		return
	}

	if err = WebHandler.ValidateRequest(request, &updateShoppingBasketItem); err != nil {
		WebHandler.BadRequest(response, request, "RequestBody", "Invalid")
		return
	}
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.UpdateShoppingBasketItem(ID.ID, updateShoppingBasketItem); err != nil {
		if errors.Is(err, ErrorOutOfStock) {
			WebHandler.BadRequest(response, request, "BadRequest", err.Error())
		} else {
			WebHandler.InternalServerError(response, request, "Error", err.Error())
		}
		return
	}
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleGetShoppingBasketByIDV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.DebugCtx(request.Context(), "HandleGetShoppingBasketByIdV1")
	var err error
	var shoppingBasket ShoppingBasket
	ID, err := WebHandler.GetID(request)
	if err != nil {
		WebHandler.BadRequest(response, request, "ID not found", err.Error())
		return
	}
	Logger.Log.DebugCtx(request.Context(), "Fetching for id: %s", ID.ID.String())
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.GetShoppingBasket(ID.ID); err != nil {
		Logger.Log.DebugCtx(request.Context(), "%s", err.Error())

		WebHandler.StatusInternalServerError(response, request)
		return
	}
	Logger.Log.DebugCtx(request.Context(), "Shopping Basket fetched: %v", shoppingBasket)
	WebHandler.WriteResponse(Http.StatusOK, response, request, shoppingBasket)
}

func NewWebHandler(shoppingBasketService ShoppingBasketService) ShoppingBasketWebHandler {
	return &shoppingBasketWebHandler{shoppingBasketService}
}
