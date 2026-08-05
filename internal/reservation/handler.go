package reservation

import (
	"net/http"

	commonHandler "github.com/danyel/ecommerce/internal/common/handler"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationHandler interface {
	CreateReservation(w http.ResponseWriter, r *http.Request)
	GetReservations(w http.ResponseWriter, r *http.Request)
}

type reservationHandler struct {
	s ReservationService
}

func (h *reservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var createReservation CreateReservation
	var reservationId ReservationId
	var err error
	if err = commonHandler.ValidateRequest[CreateReservation](r, &createReservation); err != nil {
		commonHandler.StatusBadRequest(w)
		return
	}
	if reservationId, err = h.s.CreateReservation(createReservation); err != nil {
		commonHandler.StatusInternalServerError(w)
		return
	}
	commonHandler.WriteResponse(http.StatusCreated, w, reservationId)
}

func (h *reservationHandler) GetReservations(w http.ResponseWriter, _ *http.Request) {
	commonHandler.WriteResponse(http.StatusOK, w, h.s.GetReservations())
}

func NewHandler(DB *gorm.DB) ReservationHandler {
	handler := &reservationHandler{
		NewReservationService(DB),
	}
	return handler
}
