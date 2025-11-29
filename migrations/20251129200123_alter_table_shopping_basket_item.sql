-- +goose Up
-- +goose StatementBegin
AlTER TABLE ecommerce.shopping_basket_items
    DROP COLUMN name;
AlTER TABLE ecommerce.shopping_basket_items
    DROP COLUMN image_url;
ALTER TABLE ecommerce.shopping_basket_items
    ADD COLUMN PRODUCT_ID uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
