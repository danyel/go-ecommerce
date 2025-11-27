package product

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Brand       string    `json:"brand"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	Price       int       `json:"price"`
	CategoryId  uuid.UUID `json:"category_id"`
	ImageUrl    string    `json:"image_url"`
	Stock       int       `json:"stock"`
}
