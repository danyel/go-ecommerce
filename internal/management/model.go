package management

import Uuid "github.com/google/uuid"

type CreateCategory struct {
	Name     string      `json:"name"`
	Children []Uuid.UUID `json:"children"`
}

type RemoveChild struct {
	ID      Uuid.UUID `json:"id"`
	ChildID Uuid.UUID `json:"child_id"`
}

type AddChild struct {
	ID      Uuid.UUID `json:"id"`
	ChildID Uuid.UUID `json:"child_id"`
}

type CreateCms struct {
	Code     string `json:"code"`
	Value    string `json:"value"`
	Language string `json:"language"`
}

type CmsID struct {
	ID Uuid.UUID `json:"id"`
}
