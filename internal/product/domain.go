package product

import (
	Category "github.com/danyel/ecommerce/internal/category"
	Types "github.com/danyel/ecommerce/internal/common/types"
)

type Product struct {
	ID          Types.ID          `json:"id"`
	Brand       string            `json:"brand"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Code        string            `json:"code"`
	Price       Types.Price       `json:"price"`
	Category    Category.Category `json:"category"`
	ImageURL    string            `json:"image_url"`
	Stock       int               `json:"stock"`
}
