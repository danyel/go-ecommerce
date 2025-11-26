-- +goose Up
-- +goose StatementBegin
ALTER TABLE ecommerce.products DROP COLUMN SUB_CATEGORY;
ALTER TABLE ecommerce.products DROP COLUMN CATEGORY;
ALTER TABLE ecommerce.products ADD COLUMN CATEGORY_ID UUID REFERENCES ecommerce.categories(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ecommerce.products DROP COLUMN CATEGORY;
ALTER TABLE ecommerce.products ADD COLUMN CATEGORY VARCHAR(50);
ALTER TABLE ecommerce.products ADD COLUMN SUB_CATEGORY VARCHAR(50);
-- +goose StatementEnd
