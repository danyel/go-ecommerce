package port

type EventPublisher interface {
	Publish(queue string, event any) error
}
