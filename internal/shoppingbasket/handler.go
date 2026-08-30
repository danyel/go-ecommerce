package shoppingbasket

import (
	Fmt "fmt"
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketWebHandler interface {
	HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request)
	HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetShoppingBasketByIDV1(response Http.ResponseWriter, request *Http.Request)
}

type shoppingBasketWebHandler struct {
	shoppingBasketService   ShoppingBasketService
	shoppingBasketValidator Validator
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleCreateShoppingBasketV1")
	shoppingBasket, err := shoppingBasketWebHandler.shoppingBasketService.CreateShoppingBasket()
	if err != nil {
		WebHandler.InternalServerError(response, request, WebHandler.InternalServerErrorTitle, make(map[string]any))
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, shoppingBasket.ID)
}

// HandleUpdateShoppingBasketItemV1 web handler function that will update the shopping basket
func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleUpdateShoppingBasketItemV1")
	var updateShoppingBasketItem UpdateShoppingBasketItemDTO
	var err error
	var details map[string]any
	shoppingBasket := ShoppingBasket{
		ProblemDetail: WebHandler.ProblemDetail{
			Status: 0,
		},
	}
	ID, err := WebHandler.GetID(request)
	if err != nil {
		details = make(map[string]any)
		details["id"] = Fmt.Sprintf("Could not parse id from the request: %v", err.Error())
		shoppingBasket.ProblemDetail.Status = Http.StatusNotFound
		shoppingBasket.ProblemDetail.Errors = details
	}
	if !shoppingBasket.ProblemDetail.IsError() {
		if details, err = WebHandler.ValidateRequest(request, &updateShoppingBasketItem); err != nil {
			shoppingBasket.ProblemDetail.Status = Http.StatusBadRequest
			shoppingBasket.ProblemDetail.Errors = details
		}
	}
	if !shoppingBasket.ProblemDetail.IsError() {
		problemDetail := shoppingBasketWebHandler.shoppingBasketValidator.ValidateItem(ID, updateShoppingBasketItem)
		if problemDetail.IsError() {
			shoppingBasket.ProblemDetail = problemDetail
		}
	}
	if !shoppingBasket.ProblemDetail.IsError() {
		var d ShoppingBasket
		if //goland:noinspection GoDfaErrorMayBeNotNil
		d, err = shoppingBasketWebHandler.shoppingBasketService.UpdateShoppingBasketItem(ID.ID, updateShoppingBasketItem); err != nil {
			shoppingBasket.ProblemDetail.Status = Http.StatusInternalServerError
			shoppingBasket.ProblemDetail.Errors = d.ProblemDetail.Errors
		}
	}
	if shoppingBasket.ProblemDetail.IsError() {
		WebHandler.ProblemDetailResponse(response, request, shoppingBasket.ProblemDetail)
	} else {
		WebHandler.WriteResponse(Http.StatusOK, response, request, dto(shoppingBasket))
	}
}

func dto(shoppingBasket ShoppingBasket) ShoppingBasketDTO {
	items := make([]ShoppingBasketItemDTO, len(shoppingBasket.Items))
	for i, item := range shoppingBasket.Items {
		newPrice := Types.NewPrice(float64(item.Product.Price.Inclusive*Types.Float64(item.Quantity)), "EUR")
		items[i] = ShoppingBasketItemDTO{
			Name:       item.Product.Name,
			BasePrice:  item.Product.Price.DTO(),
			TotalPrice: newPrice.DTO(),
			ProductID:  item.Product.ID,
			ImageURL:   item.Product.ImageURL,
			Quantity:   item.Quantity,
		}
	}
	totalPrice := shoppingBasket.TotalPrice()
	return ShoppingBasketDTO{
		ID:         shoppingBasket.ID,
		Items:      items,
		TotalPrice: totalPrice.DTO(),
	}
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleGetShoppingBasketByIDV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.DebugCtx(request.Context(), "HandleGetShoppingBasketByIdV1")
	var err error
	var shoppingBasket ShoppingBasket
	ID, err := WebHandler.GetID(request)
	if err != nil {
		WebHandler.BadRequest(response, request, WebHandler.IdNotFoundTitle, make(map[string]any))
		return
	}
	shoppingBasket = shoppingBasketWebHandler.shoppingBasketValidator.Validate(ID)
	Logger.Log.DebugCtx(request.Context(), "Fetching for id: %s", ID.ID.String())
	if shoppingBasket.ProblemDetail.IsError() {
		WebHandler.ProblemDetailResponse(response, request, shoppingBasket.ProblemDetail)
	}
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.GetShoppingBasket(ID.ID); err != nil {
		Logger.Log.DebugCtx(request.Context(), "%s", err.Error())
		details := make(map[string]any)
		details["id"] = Fmt.Sprintf("Could not find shopping basket with id '%s'", ID.ID.String())
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, details)
		return
	}
	Logger.Log.DebugCtx(request.Context(), "Shopping Basket fetched: %v", shoppingBasket)
	WebHandler.WriteResponse(Http.StatusOK, response, request, dto(shoppingBasket))
}

func NewWebHandler(shoppingBasketService ShoppingBasketService, shoppingBasketValidator Validator) ShoppingBasketWebHandler {
	return &shoppingBasketWebHandler{shoppingBasketService, shoppingBasketValidator}
}
