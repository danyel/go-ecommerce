package reservation

import (
	Fmt "fmt"
	Http "net/http"

	WebHandler "github.com/danyel/ecommerce/internal/common/handler"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationHandler interface {
	Create(response Http.ResponseWriter, request *Http.Request)
	FindAll(response Http.ResponseWriter, request *Http.Request)
}

type reservationHandler struct {
	reservationService ReservationService
}

func (reservationHandler *reservationHandler) Create(response Http.ResponseWriter, request *Http.Request) {
	var createReservation CreateReservation
	var ID Uuid.UUID
	var err error
	var details map[string]any
	if details, err = WebHandler.ValidateRequest[CreateReservation](request, &createReservation); err != nil {
		WebHandler.BadRequest(response, request, WebHandler.BadRequestTitle, details)
		return
	}
	reservation := Reservation{
		ShoppingBasketID: createReservation.ShoppingBasketID,
		ProductID:        createReservation.ProductID,
		Quantity:         createReservation.Quantity,
	}
	if ID, err = reservationHandler.reservationService.Create(reservation); err != nil {
		details := make(map[string]any)
		details["database"] = Fmt.Sprintf("Could not create Reservation: %s", err.Error())
		WebHandler.InternalServerError(response, request, WebHandler.InternalServerErrorTitle, details)
		return
	}
	WebHandler.WriteResponse(Http.StatusCreated, response, request, ID)
}

func (reservationHandler *reservationHandler) FindAll(response Http.ResponseWriter, request *Http.Request) {
	WebHandler.WriteResponse(Http.StatusOK, response, request, reservationHandler.reservationService.FindAll())
}

//goland:noinspection GoUnusedExportedFunction
func NewHandler(reservationService ReservationService) ReservationHandler {
	handler := &reservationHandler{
		reservationService: reservationService,
	}
	return handler
}
