package types

import (
	JSON "encoding/json"

	Uuid "github.com/google/uuid"
)

type Id struct {
	ID Uuid.UUID
}

func NewID(val Uuid.UUID) Id {
	return Id{val}
}

//goland:noinspection GoMixedReceiverTypes
func (i Id) MarshalJSON() ([]byte, error) {
	return JSON.Marshal(i.ID)
}

func (i *Id) UnmarshalJSON(data []byte) error {
	return JSON.Unmarshal(data, &i.ID)
}
