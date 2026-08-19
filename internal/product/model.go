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

type UpdateProduct struct {
	Brand       string   `json:"brand"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	CategoryID  Types.ID `json:"category_id"`
	ImageURL    string   `json:"image_url"`
	Stock       int      `json:"stock"`
}

type CreateProduct struct {
	Brand         string   `json:"brand"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Code          string   `json:"code"`
	Price         float64  `json:"price"`
	TaxPercentage int16    `json:"tax_percentage"`
	CategoryID    Types.ID `json:"category_id"`
	ImageURL      string   `json:"image_url"`
}

type Translations struct{}
