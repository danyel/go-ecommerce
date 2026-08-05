package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Id struct {
	ID uuid.UUID
}

func NewID(val uuid.UUID) Id {
	return Id{val}
}

//goland:noinspection GoMixedReceiverTypes
func (i Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ID)
}

func (i *Id) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.ID)
}
