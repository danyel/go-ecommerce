package category

import Uuid "github.com/google/uuid"

type CreateCategory struct {
	Name     string   `json:"name"`
	Children []string `json:"children"`
}

type Category struct {
	ID       Uuid.UUID  `json:"code"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}

//goland:noinspection GoNameStartsWithPackageName
type CategoryID struct {
	ID Uuid.UUID `json:"id"`
}
