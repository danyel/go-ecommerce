-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS ecommerce;
CREATE TABLE ecommerce.PRODUCTS (
    ID UUID PRIMARY KEY,
    NAME TEXT NOT NULL,
    DESCRIPTION TEXT NOT NULL,
    BRAND VARCHAR(50) NOT NULL,
    CODE VARCHAR(50) NOT NULL,
    PRICE INT default 0,
    STOCK INT default 0,
    IMAGE_URL TEXT,
    CATEGORY VARCHAR(50),
    SUB_CATEGORY VARCHAR(50),
    CREATED_AT TIMESTAMP default now(),
    UPDATED_AT TIMESTAMP default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ecommerce.PRODUCTS;
-- +goose StatementEnd
