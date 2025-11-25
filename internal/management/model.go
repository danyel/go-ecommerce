package management

import "github.com/google/uuid"

type CreateCategory struct {
	Name     string      `json:"name"`
	Children []uuid.UUID `json:"children"`
}

type RemoveChild struct {
	ID      uuid.UUID `json:"id"`
	ChildId uuid.UUID `json:"child_id"`
}

type AddChild struct {
	ID      uuid.UUID `json:"id"`
	ChildId uuid.UUID `json:"child_id"`
}

type CreateCms struct {
	Value    string `json:"value"`
	Language string `json:"language"`
}

type CmsId struct {
	ID uuid.UUID `json:"id"`
}
