-- +goose Up
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN SUB_CATEGORY;
ALTER TABLE products DROP COLUMN CATEGORY;
ALTER TABLE products ADD COLUMN CATEGORY UUID REFERENCES categories(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN CATEGORY;
ALTER TABLE products ADD COLUMN CATEGORY VARCHAR(50);
ALTER TABLE products ADD COLUMN SUB_CATEGORY VARCHAR(50);
-- +goose StatementEnd
