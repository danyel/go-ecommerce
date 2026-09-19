package handler

import (
	Fmt "fmt"
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Domain "github.com/danyel/ecommerce/internal/domain"
	Model "github.com/danyel/ecommerce/internal/model"
	Service "github.com/danyel/ecommerce/internal/service"
	Types "github.com/danyel/ecommerce/internal/types"
	Validator "github.com/danyel/ecommerce/internal/validator"
)

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketWebHandler interface {
	HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request)
	HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request)
	HandleGetShoppingBasketByIDV1(response Http.ResponseWriter, request *Http.Request)
}

type shoppingBasketWebHandler struct {
	shoppingBasketService   Service.IShoppingBasketService
	shoppingBasketValidator Validator.Validator
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleCreateShoppingBasketV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleCreateShoppingBasketV1")
	shoppingBasket, err := shoppingBasketWebHandler.shoppingBasketService.Create()
	if err != nil {
		InternalServerError(response, request, InternalServerErrorTitle, make(map[string]any))
		return
	}
	StatusCreated(response, request, shoppingBasket)
}

// HandleUpdateShoppingBasketItemV1 web handler function that will update the shopping basket
func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleUpdateShoppingBasketItemV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.Debug("HandleUpdateShoppingBasketItemV1")
	var updateShoppingBasketItem Model.UpdateShoppingBasketItemDTO
	var err error
	var details map[string]any
	shoppingBasket := Domain.ShoppingBasket{
		ProblemDetail: ProblemDetail{
			Status: 0,
		},
	}
	ID, err := GetID(request)
	if err != nil {
		details = make(map[string]any)
		details["id"] = Fmt.Sprintf("Could not parse id from the request: %v", err.Error())
		shoppingBasket.ProblemDetail.Status = Http.StatusNotFound
		shoppingBasket.ProblemDetail.Errors = details
	}
	if !shoppingBasket.ProblemDetail.IsError() {
		if details, err = ValidateRequest(request, &updateShoppingBasketItem); err != nil {
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
		if //goland:noinspection GoDfaErrorMayBeNotNil
		err = shoppingBasketWebHandler.shoppingBasketService.Update(ID.ID, Domain.ShoppingBasketItemUpdate{
			ProductID: updateShoppingBasketItem.ProductID.ID,
			Quantity:  updateShoppingBasketItem.Quantity,
		}); err != nil {
			shoppingBasket.ProblemDetail.Status = Http.StatusInternalServerError
			shoppingBasket.ProblemDetail.Title = InternalServerErrorTitle
			shoppingBasket.ProblemDetail.Details = err.Error()
		} else if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.FindByID(ID.ID); err != nil {
			shoppingBasket.ProblemDetail.Status = Http.StatusInternalServerError
			shoppingBasket.ProblemDetail.Title = InternalServerErrorTitle
			shoppingBasket.ProblemDetail.Details = err.Error()
		}
	}
	if shoppingBasket.ProblemDetail.IsError() {
		ProblemDetailResponse(response, request, shoppingBasket.ProblemDetail)
	} else {
		WriteResponse(Http.StatusOK, response, request, shoppingBasketDTO(shoppingBasket))
	}
}

func shoppingBasketDTO(shoppingBasket Domain.ShoppingBasket) Model.ShoppingBasketDTO {
	items := make([]Model.ShoppingBasketItemDTO, len(shoppingBasket.Items))
	for i, item := range shoppingBasket.Items {
		newPrice := Types.NewPrice(float64(item.Product.Price.Inclusive*Types.Float64(item.Quantity)), "EUR")
		items[i] = Model.ShoppingBasketItemDTO{
			Name:       item.Product.Name,
			BasePrice:  item.Product.Price.DTO(),
			TotalPrice: newPrice.DTO(),
			ProductID:  item.Product.ID,
			ImageURL:   item.Product.ImageURL,
			Quantity:   item.Quantity,
		}
	}
	totalPrice := shoppingBasket.TotalPrice()
	return Model.ShoppingBasketDTO{
		ID:         shoppingBasket.ID,
		Items:      items,
		TotalPrice: totalPrice.DTO(),
	}
}

func (shoppingBasketWebHandler *shoppingBasketWebHandler) HandleGetShoppingBasketByIDV1(response Http.ResponseWriter, request *Http.Request) {
	Logger.Log.DebugCtx(request.Context(), "HandleGetShoppingBasketByIdV1")
	var err error
	var shoppingBasket Domain.ShoppingBasket
	ID, err := GetID(request)
	Logger.Log.DebugCtx(request.Context(), "HandleGetShoppingBasketByIdV1: %s %v", ID, err)
	if err != nil {
		BadRequest(response, request, IdNotFoundTitle, make(map[string]any))
		return
	}
	shoppingBasket = shoppingBasketWebHandler.shoppingBasketValidator.Validate(ID)
	Logger.Log.DebugCtx(request.Context(), "Fetching for id: %s", ID.ID.String())
	if shoppingBasket.ProblemDetail.IsError() {
		ProblemDetailResponse(response, request, shoppingBasket.ProblemDetail)
		return
	}
	if shoppingBasket, err = shoppingBasketWebHandler.shoppingBasketService.FindByID(ID.ID); err != nil {
		Logger.Log.DebugCtx(request.Context(), "%s", err.Error())
		details := make(map[string]any)
		details["id"] = Fmt.Sprintf("Could not find shopping basket with id '%s'", ID.ID.String())
		BadRequest(response, request, BadRequestTitle, details)
		return
	}
	Logger.Log.DebugCtx(request.Context(), "Shopping Basket fetched: %v", shoppingBasket)
	WriteResponse(Http.StatusOK, response, request, shoppingBasketDTO(shoppingBasket))
}

func NewShoppingBasketWebHandler(shoppingBasketService Service.IShoppingBasketService, shoppingBasketValidator Validator.Validator) ShoppingBasketWebHandler {
	return &shoppingBasketWebHandler{shoppingBasketService, shoppingBasketValidator}
}
