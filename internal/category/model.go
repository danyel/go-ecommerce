package category

import "github.com/google/uuid"

type CreateCategory struct {
	Name     string      `json:"name"`
	Children []uuid.UUID `json:"children"`
}

type Category struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}

type RemoveChild struct {
	ID      uuid.UUID `json:"id"`
	ChildId uuid.UUID `json:"child_id"`
}

type AddChild struct {
	ID      uuid.UUID `json:"id"`
	ChildId uuid.UUID `json:"child_id"`
}

type CategoryId struct {
	ID uuid.UUID `json:"id"`
}
