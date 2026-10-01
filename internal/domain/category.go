package domain

import (
	Types "github.com/danyel/ecommerce/internal/types"
)

type Category struct {
	ID       Types.ID   `json:"id"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}
