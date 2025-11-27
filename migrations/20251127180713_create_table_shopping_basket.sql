-- +goose Up
-- +goose StatementBegin
CREATE TABLE ecommerce.shopping_basket(ID uuid PRIMARY KEY, created_at TIMESTAMP default now(), updated_at TIMESTAMP DEFAULT now());
CREATE TABLE ecommerce.shopping_basket_items(shopping_basket_id UUID references ecommerce.shopping_basket(ID), product_id UUID references ecommerce.products(id));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ecommerce.shopping_basket;
DROP TABLE ecommerce.shopping_basket_items;
-- +goose StatementEnd
