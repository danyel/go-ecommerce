package reservation

import (
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationService interface {
	GetReservations() []Reservation
	GetReservation(reservationID Types.ID) (Reservation, error)
	CreateReservation(createReservation CreateReservation) (Types.ID, error)
}

type reservationService struct {
	reservationRepository Repository.CrudRepository[ReservationModel]
}

func (s *reservationService) GetReservations() []Reservation {
	reservationModels := s.reservationRepository.FindAll(Repository.SearchCriteria{Preloads: []string{"Children"}})
	return mapReservations(reservationModels)
}

func (s *reservationService) GetReservation(reservationID Types.ID) (Reservation, error) {
	var reservation Reservation
	reservationModel, err := s.reservationRepository.FindById(reservationID.ID)
	if err != nil {
		return reservation, err
	}
	return mapReservation(reservationModel), err
}

func (s *reservationService) CreateReservation(createReservation CreateReservation) (Types.ID, error) {
	var err error
	reservation := &ReservationModel{
		ShoppingBasketID: createReservation.ShoppingBasketID.ID,
		ProductID:        createReservation.ProductID.ID,
		Quantity:         createReservation.Quantity,
	}

	if err := s.reservationRepository.Create(reservation); err != nil {
		return Types.NewID(reservation.ShoppingBasketID), err
	}
	return Types.NewID(reservation.ShoppingBasketID), err
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
		ShoppingBasketID: Types.NewID(reservationModel.ShoppingBasketID),
		ProductID:        Types.NewID(reservationModel.ProductID),
		Quantity:         reservationModel.Quantity,
	}
}

func NewService(reservationRepository Repository.CrudRepository[ReservationModel]) ReservationService {
	return &reservationService{
		reservationRepository: reservationRepository,
	}
}
