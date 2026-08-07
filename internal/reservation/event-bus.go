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

//goland:noinspection GoUnusedGlobalVariable,GoNameStartsWithPackageName
var ReservationCreated = Broker.QueueConfig{
	Topic: ExchangeReservation,
	Queue: AddReservationQueue,
}

//goland:noinspection GoNameStartsWithPackageName
type ReservationCreatedEvent struct {
	Id string `json:"id"`
}

//goland:noinspection GoUnusedExportedFunction
func HandleReservationCreated2(body []byte) error {
	var event ReservationCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Println(event.Id)
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func HandleReservationCreated(body []byte) error {
	var event ReservationCreatedEvent
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Log.Println(event.Id)
	return nil
}
