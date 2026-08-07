package category

import Uuid "github.com/google/uuid"

type CategoryAggregate struct {
	AggregateIdentifier Uuid.UUID
	Name                string
	Children            []Uuid.UUID
}
