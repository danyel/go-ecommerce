package reservation

import (
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ReservationService interface {
	FindAll() []Reservation
	Find(reservationID Uuid.UUID) (Reservation, error)
	Create(reservation Reservation) (Uuid.UUID, error)
	Update(shoppingBasketID Uuid.UUID, productID Uuid.UUID, quantity int) error
}

type reservationService struct {
	reservationRepository Repository.CrudRepository[ReservationModel]
}

func (reservationService *reservationService) FindAll() []Reservation {
	reservationModels := reservationService.reservationRepository.FindAll(Repository.SearchCriteria{Preloads: []string{"Children"}})
	return mapReservations(reservationModels)
}

func (reservationService *reservationService) Find(reservationID Uuid.UUID) (Reservation, error) {
	var reservation Reservation
	reservationModel, err := reservationService.reservationRepository.FindById(reservationID)
	if err != nil {
		return reservation, err
	}
	return mapReservation(reservationModel), err
}

func (reservationService *reservationService) Create(reservation Reservation) (Uuid.UUID, error) {
	var err error
	reservationModel := &ReservationModel{
		ShoppingBasketID: reservation.ShoppingBasketID.ID,
		ProductID:        reservation.ProductID.ID,
		Quantity:         reservation.Quantity,
	}

	if err := reservationService.reservationRepository.Create(reservationModel); err != nil {
		return reservationModel.ShoppingBasketID, err
	}
	return reservationModel.ShoppingBasketID, err
}

func (reservationService *reservationService) Update(shoppingBasketID Uuid.UUID, productID Uuid.UUID, quantity int) error {
	clause := make([]any, 2)
	clause[0] = shoppingBasketID.String()
	clause[0] = productID.String()
	reservations := reservationService.reservationRepository.FindAll(Repository.SearchCriteria{
		WhereClause: Repository.WhereClause{
			Query:  "shopping_basket_id = ? AND product_id = ?",
			Params: clause,
		},
	})

	var reservation = &ReservationModel{
		ShoppingBasketID: shoppingBasketID,
		ProductID:        productID,
		Quantity:         quantity,
	}

	if len(reservations) == 0 {
		err := reservationService.reservationRepository.Create(reservation)
		if err != nil {
			return err
		}
	} else {
		reservation = reservations[0]
		reservation.Quantity = quantity
	}

	return reservationService.reservationRepository.Update(reservation)
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
