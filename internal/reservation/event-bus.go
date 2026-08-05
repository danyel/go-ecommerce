package reservation

import (
	"encoding/json"
	"log"

	"github.com/danyel/ecommerce/cmd/broker"
)

const (
	ExchangeReservation = "reservation.topic"
	AddReservationQueue = "reservations.add_reservation"
)

var ReservationCreated = broker.QueueConfig{
	Topic: ExchangeReservation,
	Queue: AddReservationQueue,
}

type ReservationCreatedEvent struct {
	Id string `json:"id"`
}

func HandleReservationCreated2(body []byte) error {
	var event ReservationCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	log.Println(event.Id)
	return nil
}
func HandleReservationCreated(body []byte) error {
	var event ReservationCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	log.Println(event.Id)
	return nil
}
