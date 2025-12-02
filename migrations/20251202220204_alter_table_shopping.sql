-- +goose Up
-- +goose StatementBegin
ALTER TABLE ecommerce.shopping_basket_items RENAME COLUMN amount to quantity;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
