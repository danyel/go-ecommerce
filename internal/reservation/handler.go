package reservation

import (
	Http "net/http"

	CommonHandler "github.com/danyel/ecommerce/internal/common/handler"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationHandler interface {
	CreateReservation(w Http.ResponseWriter, r *Http.Request)
	GetReservations(w Http.ResponseWriter, r *Http.Request)
}

type reservationHandler struct {
	s ReservationService
}

func (h *reservationHandler) CreateReservation(w Http.ResponseWriter, r *Http.Request) {
	var createReservation CreateReservation
	var reservationID Types.ID
	var err error
	if err = CommonHandler.ValidateRequest[CreateReservation](r, &createReservation); err != nil {
		CommonHandler.StatusBadRequest(w, r)
		return
	}
	if reservationID, err = h.s.CreateReservation(createReservation); err != nil {
		CommonHandler.StatusInternalServerError(w, r)
		return
	}
	CommonHandler.WriteResponse(Http.StatusCreated, w, r, reservationID)
}

func (h *reservationHandler) GetReservations(w Http.ResponseWriter, r *Http.Request) {
	CommonHandler.WriteResponse(Http.StatusOK, w, r, h.s.GetReservations())
}

// NewHandler adding to router (todo)
//
//goland:noinspection GoUnusedExportedFunction
func NewHandler(DB *Database.DB) ReservationHandler {
	handler := &reservationHandler{
		NewReservationService(DB),
	}
	return handler
}
