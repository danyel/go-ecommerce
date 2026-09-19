package model

import (
	Types "github.com/danyel/ecommerce/internal/types"
)

type CreateCategory struct {
	Name     string   `json:"name" validate:"required,min=3"`
	Children []string `json:"children"`
}

type CategoryDTO struct {
	ID       Types.ID      `json:"id"`
	Name     string        `json:"name"`
	Children []CategoryDTO `json:"children,omitempty"`
}

//goland:noinspection GoNameStartsWithPackageName
type CategoryID struct {
	ID Types.ID `json:"id"`
}
