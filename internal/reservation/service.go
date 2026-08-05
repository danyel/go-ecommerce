package reservation

import "C"
import (
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationService interface {
	GetReservations() []Reservation
	GetReservation(reservationID uuid.UUID) (Reservation, error)
	CreateReservation(createReservation CreateReservation) (ReservationId, error)
}

type reservationService struct {
	reservationRepository commonRepository.CrudRepository[ReservationModel]
}

func (s *reservationService) GetReservations() []Reservation {
	reservationModels := s.reservationRepository.FindAll(commonRepository.SearchCriteria{Preloads: []string{"Children"}})
	return mapReservations(reservationModels)
}

func (s *reservationService) GetReservation(reservationID uuid.UUID) (Reservation, error) {
	var reservation Reservation
	reservationModel, err := s.reservationRepository.FindById(reservationID)
	if err != nil {
		return reservation, err
	}
	return mapReservation(reservationModel), err
}

func (s *reservationService) CreateReservation(createReservation CreateReservation) (ReservationId, error) {
	var err error
	var reservationId ReservationId
	reservation := &ReservationModel{
		ShoppingBasketId: createReservation.ShoppingBasketId,
		ProductId:        createReservation.ProductId,
		Quantity:         createReservation.Quantity,
	}

	if err := s.reservationRepository.Create(reservation); err != nil {
		return reservationId, err
	}
	reservationId.ID = reservation.ShoppingBasketId
	return reservationId, err
}

func mapReservations(models []*ReservationModel) []Reservation {
	reservations := make([]Reservation, len(models))

	for i, m := range models {
		reservations[i] = mapReservation(m)
	}

	return reservations
}

func mapReservation(reservationModel *ReservationModel) Reservation {
	return Reservation{
		ShoppingBasketId: reservationModel.ShoppingBasketId,
		ProductId:        reservationModel.ProductId,
		Quantity:         reservationModel.Quantity,
	}
}

func NewReservationService(DB *gorm.DB) ReservationService {
	return &reservationService{
		reservationRepository: commonRepository.NewCrudRepository[ReservationModel](DB),
	}
}
