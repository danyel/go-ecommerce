package product

import (
	JSON "encoding/json"

	Broker "github.com/danyel/ecommerce/cmd/broker"
	Logger "github.com/danyel/ecommerce/cmd/logger"
)

const (
	ExchangeReservation = "product.topic"
	AddReservationQueue = "products.update_stock"
)

type CommandHandler interface {
	handleReservationCreate(body []byte) error
}

//goland:noinspection GoUnusedGlobalVariable,GoNameStartsWithPackageName
var UpdateProductStock = Broker.QueueConfig{
	Topic: ExchangeReservation,
	Queue: AddReservationQueue,
}

func (commandHandler *commandHandler) handleUpdateStock(body []byte) error {
	var event UpdateProductStockCommand
	if err := JSON.Unmarshal(body, &event); err != nil {
		return err
	}
	Logger.Log.Debug("Entering: %s with product id: %s", AddReservationQueue, event.ProductID)
	product, err := commandHandler.productService.FindById(event.ProductID)
	if err != nil {
		return err
	}
	if event.Quantity > 0 {
		product.Stock -= event.Quantity
	} else if event.Quantity < 0 {
		product.Stock += -event.Quantity
	}
	return commandHandler.productService.Update(product)
}

type commandHandler struct {
	productService ProductService
}

func RegisterConsumer(productService ProductService, messageBroker *Broker.MessageBroker) {
	commandHandlerInstance := &commandHandler{productService}
	messageBroker.RegisterConsumer(UpdateProductStock, commandHandlerInstance.handleUpdateStock)
}
