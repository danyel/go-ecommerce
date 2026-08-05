package product

import (
	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/common/types"
)

type Product struct {
	ID          types.Id          `json:"id"`
	Brand       string            `json:"brand"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Code        string            `json:"code"`
	Price       float64           `json:"price"`
	Category    category.Category `json:"category"`
	ImageUrl    string            `json:"image_url"`
	Stock       int               `json:"stock"`
}
