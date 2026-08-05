package reservation

import "C"
import (
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"github.com/danyel/ecommerce/internal/common/types"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationService interface {
	GetReservations() []Reservation
	GetReservation(reservationID types.Id) (Reservation, error)
	CreateReservation(createReservation CreateReservation) (types.Id, error)
}

type reservationService struct {
	reservationRepository commonRepository.CrudRepository[ReservationModel]
}

func (s *reservationService) GetReservations() []Reservation {
	reservationModels := s.reservationRepository.FindAll(commonRepository.SearchCriteria{Preloads: []string{"Children"}})
	return mapReservations(reservationModels)
}

func (s *reservationService) GetReservation(reservationID types.Id) (Reservation, error) {
	var reservation Reservation
	reservationModel, err := s.reservationRepository.FindById(reservationID.ID)
	if err != nil {
		return reservation, err
	}
	return mapReservation(reservationModel), err
}

func (s *reservationService) CreateReservation(createReservation CreateReservation) (types.Id, error) {
	var err error
	reservation := &ReservationModel{
		ShoppingBasketId: createReservation.ShoppingBasketId.ID,
		ProductId:        createReservation.ProductId.ID,
		Quantity:         createReservation.Quantity,
	}

	if err := s.reservationRepository.Create(reservation); err != nil {
		return types.NewID(reservation.ShoppingBasketId), err
	}
	return types.NewID(reservation.ShoppingBasketId), err
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
		ShoppingBasketId: types.NewID(reservationModel.ShoppingBasketId),
		ProductId:        types.NewID(reservationModel.ProductId),
		Quantity:         reservationModel.Quantity,
	}
}

func NewReservationService(DB *gorm.DB) ReservationService {
	return &reservationService{
		reservationRepository: commonRepository.NewCrudRepository[ReservationModel](DB),
	}
}
