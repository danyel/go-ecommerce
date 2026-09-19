// Package domain: domain layer
package domain

import (
	Types "github.com/danyel/ecommerce/internal/types"
)

type Product struct {
	ID          Types.ID    `json:"id"`
	Brand       string      `json:"brand"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Code        string      `json:"code"`
	Price       Types.Price `json:"price"`
	Category    Category    `json:"category"`
	ImageURL    string      `json:"image_url"`
	Stock       int         `json:"stock"`
}
