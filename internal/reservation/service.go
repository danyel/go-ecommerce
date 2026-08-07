package reservation

import (
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Database "gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationService interface {
	GetReservations() []Reservation
	GetReservation(reservationID Types.Id) (Reservation, error)
	CreateReservation(createReservation CreateReservation) (Types.Id, error)
}

type reservationService struct {
	reservationRepository CommonRepository.CrudRepository[ReservationModel]
}

func (s *reservationService) GetReservations() []Reservation {
	reservationModels := s.reservationRepository.FindAll(CommonRepository.SearchCriteria{Preloads: []string{"Children"}})
	return mapReservations(reservationModels)
}

func (s *reservationService) GetReservation(reservationID Types.Id) (Reservation, error) {
	var reservation Reservation
	reservationModel, err := s.reservationRepository.FindById(reservationID.ID)
	if err != nil {
		return reservation, err
	}
	return mapReservation(reservationModel), err
}

func (s *reservationService) CreateReservation(createReservation CreateReservation) (Types.Id, error) {
	var err error
	reservation := &ReservationModel{
		ShoppingBasketId: createReservation.ShoppingBasketId.ID,
		ProductId:        createReservation.ProductId.ID,
		Quantity:         createReservation.Quantity,
	}

	if err := s.reservationRepository.Create(reservation); err != nil {
		return Types.NewID(reservation.ShoppingBasketId), err
	}
	return Types.NewID(reservation.ShoppingBasketId), err
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
		ShoppingBasketId: Types.NewID(reservationModel.ShoppingBasketId),
		ProductId:        Types.NewID(reservationModel.ProductId),
		Quantity:         reservationModel.Quantity,
	}
}

func NewReservationService(DB *Database.DB) ReservationService {
	return &reservationService{
		reservationRepository: CommonRepository.NewCrudRepository[ReservationModel](DB),
	}
}
