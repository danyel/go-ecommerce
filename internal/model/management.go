package model

import Uuid "github.com/google/uuid"

type RemoveChild struct {
	ID      Uuid.UUID `json:"id"`
	ChildID Uuid.UUID `json:"child_id"`
}

type AddChild struct {
	ID      Uuid.UUID `json:"id"`
	ChildID Uuid.UUID `json:"child_id"`
}
