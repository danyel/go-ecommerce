-- +goose Up
-- +goose StatementBegin
CREATE TABLE ecommerce.CATEGORIES
(
    ID         UUID PRIMARY KEY,
    NAME       TEXT NOT NULL,
    CREATED_AT TIMESTAMP default now(),
    UPDATED_AT TIMESTAMP default now()
);

CREATE TABLE IF NOT EXISTS ecommerce.category_children (
    parent_id UUID NOT NULL REFERENCES ecommerce.categories(id) ON DELETE CASCADE,
    child_id  UUID NOT NULL REFERENCES ecommerce.categories(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_id, child_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ecommerce.CATEGORY_CHILDREN;
DROP TABLE ecommerce.CATEGORIES;
-- +goose StatementEnd
