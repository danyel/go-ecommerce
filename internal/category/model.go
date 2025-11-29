package category

import "github.com/google/uuid"

type CreateCategory struct {
	Name     string   `json:"name"`
	Children []string `json:"children"`
}

type Category struct {
	ID       uuid.UUID  `json:"code"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}

//goland:noinspection GoNameStartsWithPackageName
type CategoryId struct {
	ID uuid.UUID `json:"id"`
}
