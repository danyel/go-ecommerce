package reservation

import (
	Fmt "fmt"
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
	var details map[string]any
	if details, err = WebHandler.ValidateRequest[CreateReservation](request, &createReservation); err != nil {
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, details)
		return
	}
	if ID, err = reservationHandler.reservationService.CreateReservation(createReservation); err != nil {
		details := make(map[string]any)
		details["database"] = Fmt.Sprintf("Could not create Reservation: %s", err.Error())
		WebHandler.InternalServerError(response, request, WebHandler.InternalServerErrorTitle, details)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, ID)
}

func (reservationHandler *reservationHandler) GetReservations(response Http.ResponseWriter, request *Http.Request) {
	WebHandler.WriteResponse(Http.StatusOK, response, request, reservationHandler.reservationService.GetReservations())
}

//goland:noinspection GoUnusedExportedFunction
func NewHandler(reservationService ReservationService) ReservationHandler {
	handler := &reservationHandler{
		reservationService: reservationService,
	}
	return handler
}
