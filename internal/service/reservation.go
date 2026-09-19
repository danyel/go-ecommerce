package service

import (
	Model "github.com/danyel/ecommerce/internal/model"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Types "github.com/danyel/ecommerce/internal/types"
	Uuid "github.com/google/uuid"
)

type IReservationService interface {
	FindAll() []Model.Reservation
	Find(reservationID Uuid.UUID) (Model.Reservation, error)
	Create(reservation Model.Reservation) (Uuid.UUID, error)
	Update(shoppingBasketID Uuid.UUID, productID Uuid.UUID, quantity int) error
}

type reservationService struct {
	reservationRepository Persistence.CrudRepository[Persistence.ReservationModel]
}

func (reservationService *reservationService) FindAll() []Model.Reservation {
	reservationModels := reservationService.reservationRepository.FindAll(Persistence.SearchCriteria{Preloads: []string{"Children"}})
	return mapReservations(reservationModels)
}

func (reservationService *reservationService) Find(reservationID Uuid.UUID) (Model.Reservation, error) {
	var reservation Model.Reservation
	reservationModel, err := reservationService.reservationRepository.FindByID(reservationID)
	if err != nil {
		return reservation, err
	}
	return mapReservation(reservationModel), err
}

func (reservationService *reservationService) Create(reservation Model.Reservation) (Uuid.UUID, error) {
	var err error
	reservationModel := &Persistence.ReservationModel{
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
	reservations := reservationService.reservationRepository.FindAll(Persistence.SearchCriteria{
		WhereClause: Persistence.WhereClause{
			Query:  "shopping_basket_id = ? AND product_id = ?",
			Params: clause,
		},
	})

	var reservation = &Persistence.ReservationModel{
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

func mapReservations(models []*Persistence.ReservationModel) []Model.Reservation {
	reservations := make([]Model.Reservation, len(models))

	for i, m := range models {
		reservations[i] = mapReservation(m)
	}

	return reservations
}

func mapReservation(reservationModel *Persistence.ReservationModel) Model.Reservation {
	return Model.Reservation{
		ShoppingBasketID: Types.NewID(reservationModel.ShoppingBasketID),
		ProductID:        Types.NewID(reservationModel.ProductID),
		Quantity:         reservationModel.Quantity,
	}
}

func ReservationService(reservationRepository Persistence.CrudRepository[Persistence.ReservationModel]) IReservationService {
	return &reservationService{
		reservationRepository: reservationRepository,
	}
}
