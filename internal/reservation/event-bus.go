package reservation

import (
	JSON "encoding/json"
	Log "log"

	Broker "github.com/danyel/ecommerce/cmd/broker"
)

const (
	ExchangeReservation = "reservation.topic"
	AddReservationQueue = "reservations.add_reservation"
)

var ReservationCreated = Broker.QueueConfig{
	Topic: ExchangeReservation,
	Queue: AddReservationQueue,
}

type ReservationCreatedEvent struct {
	Id string `json:"id"`
}

func HandleReservationCreated2(body []byte) error {
	var event ReservationCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Println(event.Id)
	return nil
}
func HandleReservationCreated(body []byte) error {
	var event ReservationCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Println(event.Id)
	return nil
}
