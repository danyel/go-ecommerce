package category

import "github.com/google/uuid"

type CategoryAggregate struct {
	AggregateIdentifier uuid.UUID
	Name                string
	Children            []uuid.UUID
}
