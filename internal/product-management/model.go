package product_management

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
	Category    Category  `json:"category"`
	ImageUrl    string    `json:"image_url"`
	Stock       int       `json:"stock"`
	//Translations []Translations `json:"translations"`
}

type UpdateProduct struct {
	Brand       string    `json:"brand"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CategoryId  uuid.UUID `json:"category_id"`
	ImageUrl    string    `json:"image_url"`
	Stock       int       `json:"stock"`
}

type CreateProduct struct {
	Brand       string    `json:"brand"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	Price       int       `json:"price"`
	CategoryId  uuid.UUID `json:"category_id"`
	ImageUrl    string    `json:"image_url"`
}

type ProductId struct {
	ID uuid.UUID `json:"id"`
}

type Category struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}

type Translations struct {
}
