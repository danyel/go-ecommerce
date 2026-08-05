package port

type EventPublisher interface {
	Publish(queue string, v interface{}) error
}
