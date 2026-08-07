package category

import Uuid "github.com/google/uuid"

//goland:noinspection GoNameStartsWithPackageName
type CategoryAggregate struct {
	AggregateIdentifier Uuid.UUID
	Name                string
	Children            []Uuid.UUID
}
