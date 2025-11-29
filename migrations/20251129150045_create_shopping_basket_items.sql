-- +goose Up
-- +goose StatementBegin
DROP TABLE ecommerce.shopping_basket_items;
CREATE TABLE ecommerce.shopping_basket_items(id uuid primary key, shopping_basket_id UUID references ecommerce.shopping_basket(ID), name text, amount integer default 0, image_url text, price integer default 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
