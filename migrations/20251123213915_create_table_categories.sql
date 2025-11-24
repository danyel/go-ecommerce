-- +goose Up
-- +goose StatementBegin
CREATE TABLE CATEGORIES
(
    ID         UUID PRIMARY KEY,
    NAME       TEXT NOT NULL,
    CREATED_AT TIMESTAMP default now(),
    UPDATED_AT TIMESTAMP default now()
);

CREATE TABLE IF NOT EXISTS category_children (
    parent_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    child_id  UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_id, child_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE CATEGORY_CHILDREN;
DROP TABLE CATEGORIES;
-- +goose StatementEnd
