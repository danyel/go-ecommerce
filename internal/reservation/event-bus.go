package reservation

import (
	JSON "encoding/json"

	Broker "github.com/danyel/ecommerce/cmd/broker"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Product "github.com/danyel/ecommerce/internal/product"
)

const (
	ExchangeReservation = "reservation.topic"
	AddReservationQueue = "reservations.add_reservation"
)

type ReservationEventHandler interface {
	handleReservationCreate(body []byte) error
}

//goland:noinspection GoUnusedGlobalVariable,GoNameStartsWithPackageName
var ReservationCreated = Broker.QueueConfig{
	Topic: ExchangeReservation,
	Queue: AddReservationQueue,
}

//goland:noinspection GoNameStartsWithPackageName
func (h *reservationEvents) handleReservationCreated(body []byte) error {
	var event ReservationCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Logger.Log.Debug("%s", event.ID)

	// create reservation if it does not exist

	// alter the quantity of the product

	return nil
}

type reservationEvents struct {
	s ReservationService
	p Product.ProductService
}

func RegisterReservationEvents(reservationService ReservationService, productService Product.ProductService, b *Broker.MessageBroker) {
	h := &reservationEvents{reservationService, productService}
	b.RegisterConsumer(ReservationCreated, h.handleReservationCreated)
}
