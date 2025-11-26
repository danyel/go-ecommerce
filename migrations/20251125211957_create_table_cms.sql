-- +goose Up
-- +goose StatementBegin
CREATE TABLE ecommerce.CMS
(
    ID         UUID PRIMARY KEY,
    CODE       TEXT       NOT NULL,
    VALUE      TEXT       NOT NULL,
    LANGUAGE   VARCHAR(5) NOT NULL,
    CREATED_AT TIMESTAMP DEFAULT now(),
    UPDATED_AT TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ecommerce.CMS;
-- +goose StatementEnd
