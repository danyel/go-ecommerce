package category

import (
	Types "github.com/danyel/ecommerce/internal/common/types"
)

type CreateCategory struct {
	Name     string   `json:"name" validate:"required,min=3"`
	Children []string `json:"children"`
}

type Category struct {
	ID       Types.ID   `json:"id"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}

//goland:noinspection GoNameStartsWithPackageName
type CategoryID struct {
	ID Types.ID `json:"id"`
}
