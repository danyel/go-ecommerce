package category

import "github.com/google/uuid"

type CreateCategory struct {
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Children []uuid.UUID `json:"children"`
}

type Category struct {
	Code     string     `json:"code"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}

//goland:noinspection GoNameStartsWithPackageName
type CategoryId struct {
	ID uuid.UUID `json:"id"`
}
