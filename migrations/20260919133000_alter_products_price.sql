-- +goose Up
-- +goose StatementBegin
ALTER TABLE ecommerce.products
    ALTER COLUMN price TYPE NUMERIC(10, 2)
    USING price::numeric;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ecommerce.products
    ALTER COLUMN price TYPE INT
    USING round(price)::int;
-- +goose StatementEnd
