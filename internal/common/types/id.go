package types

import (
	JSON "encoding/json"

	Uuid "github.com/google/uuid"
)

type ID struct {
	ID Uuid.UUID
}

func NewID(val Uuid.UUID) ID {
	return ID{val}
}

//goland:noinspection GoMixedReceiverTypes
func (i ID) MarshalJSON() ([]byte, error) {
	return JSON.Marshal(i.ID)
}

func (i *ID) UnmarshalJSON(data []byte) error {
	return JSON.Unmarshal(data, &i.ID)
}
