package handler

import (
	Fmt "fmt"
	Http "net/http"

	Model "github.com/danyel/ecommerce/internal/model"
	Service "github.com/danyel/ecommerce/internal/service"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationHandler interface {
	Create(response Http.ResponseWriter, request *Http.Request)
	FindAll(response Http.ResponseWriter, request *Http.Request)
}

type reservationHandler struct {
	reservationService Service.IReservationService
}

func (reservationHandler *reservationHandler) Create(response Http.ResponseWriter, request *Http.Request) {
	var createReservation Model.CreateReservation
	var ID Uuid.UUID
	var err error
	var details map[string]any
	if details, err = ValidateRequest[Model.CreateReservation](request, &createReservation); err != nil {
		BadRequest(response, request, BadRequestTitle, details)
		return
	}
	reservation := Model.Reservation{
		ShoppingBasketID: createReservation.ShoppingBasketID,
		ProductID:        createReservation.ProductID,
		Quantity:         createReservation.Quantity,
	}
	if ID, err = reservationHandler.reservationService.Create(reservation); err != nil {
		details := make(map[string]any)
		details["database"] = Fmt.Sprintf("Could not create Reservation: %s", err.Error())
		InternalServerError(response, request, InternalServerErrorTitle, details)
		return
	}
	WriteResponse(Http.StatusCreated, response, request, ID)
}

func (reservationHandler *reservationHandler) FindAll(response Http.ResponseWriter, request *Http.Request) {
	WriteResponse(Http.StatusOK, response, request, reservationHandler.reservationService.FindAll())
}

//goland:noinspection GoUnusedExportedFunction
func NewHandler(reservationService Service.IReservationService) ReservationHandler {
	handler := &reservationHandler{
		reservationService: reservationService,
	}
	return handler
}
