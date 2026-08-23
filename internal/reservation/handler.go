package reservation

import (
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationHandler interface {
	CreateReservation(response Http.ResponseWriter, request *Http.Request)
	GetReservations(response Http.ResponseWriter, request *Http.Request)
}

type reservationHandler struct {
	reservationService ReservationService
}

func (reservationHandler *reservationHandler) CreateReservation(response Http.ResponseWriter, request *Http.Request) {
	var createReservation CreateReservation
	var ID Types.ID
	var err error
	if err = WebHandler.ValidateRequest[CreateReservation](request, &createReservation); err != nil {
		WebHandler.StatusBadRequest(response, request)
		return
	}
	if ID, err = reservationHandler.reservationService.CreateReservation(createReservation); err != nil {
		WebHandler.StatusInternalServerError(response, request)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, ID)
}

func (reservationHandler *reservationHandler) GetReservations(response Http.ResponseWriter, request *Http.Request) {
	WebHandler.WriteResponse(Http.StatusOK, response, request, reservationHandler.reservationService.GetReservations())
}

// NewHandler adding to router (todo)
//
//goland:noinspection GoUnusedExportedFunction
func NewHandler(reservationService ReservationService) ReservationHandler {
	handler := &reservationHandler{
		reservationService: reservationService,
	}
	return handler
}
