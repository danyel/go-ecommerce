-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE SCHEMA IF NOT EXISTS ecommerce;
SET search_path TO ecommerce, public;

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email TEXT,
                       display_name TEXT NOT NULL DEFAULT '',
                       role TEXT NOT NULL DEFAULT 'USER' CHECK (role IN ('USER', 'ADMIN', 'CATALOG_MANAGER')),
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_unique ON users (lower(email)) WHERE email IS NOT NULL;

CREATE TABLE identities (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            provider TEXT NOT NULL,
                            subject TEXT NOT NULL,
                            email TEXT,
                            created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                            updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                            UNIQUE (provider, subject)
);
CREATE INDEX identities_user_id_idx ON identities(user_id);

CREATE TABLE sessions (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          token_hash TEXT NOT NULL UNIQUE,
                          expires_at TIMESTAMPTZ NOT NULL,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                          revoked_at TIMESTAMPTZ
);
CREATE INDEX sessions_user_id_idx ON sessions(user_id);

CREATE TABLE categories (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            parent_id UUID REFERENCES categories(id) ON DELETE RESTRICT,
                            name TEXT NOT NULL,
                            slug TEXT NOT NULL UNIQUE,
                            created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                            updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX categories_parent_id_idx ON categories(parent_id);

CREATE TABLE products (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
                          brand TEXT NOT NULL,
                          name TEXT NOT NULL,
                          description TEXT NOT NULL DEFAULT '',
                          code TEXT NOT NULL UNIQUE,
                          price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
                          stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
                          image_url TEXT NOT NULL DEFAULT '',
                          metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX products_category_id_idx ON products(category_id);
CREATE INDEX products_created_at_idx ON products(created_at, id);

CREATE TABLE cms (
                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                     code TEXT NOT NULL,
                     value TEXT NOT NULL,
                     language VARCHAR(10) NOT NULL,
                     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                     updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                     UNIQUE (code, language)
);

CREATE TABLE shopping_baskets (
                                  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
                                  status TEXT NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CHECKED_OUT', 'ABANDONED')),
                                  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX shopping_baskets_user_status_idx ON shopping_baskets(user_id, status);

CREATE TABLE shopping_basket_items (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       shopping_basket_id UUID NOT NULL REFERENCES shopping_baskets(id) ON DELETE CASCADE,
                                       product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
                                       unit_price NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
                                       quantity INTEGER NOT NULL CHECK (quantity > 0),
                                       max_quantity INTEGER NOT NULL DEFAULT 0 CHECK (max_quantity >= 0),
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                       UNIQUE (shopping_basket_id, product_id)
);
CREATE INDEX shopping_basket_items_product_id_idx ON shopping_basket_items(product_id);

CREATE TABLE reservations (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              shopping_basket_id UUID NOT NULL REFERENCES shopping_baskets(id) ON DELETE CASCADE,
                              product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
                              quantity INTEGER NOT NULL CHECK (quantity > 0),
                              created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                              updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                              UNIQUE (shopping_basket_id, product_id)
);

CREATE TABLE outbox_events (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               aggregate_type TEXT NOT NULL,
                               aggregate_id UUID NOT NULL,
                               event_type TEXT NOT NULL,
                               payload JSONB NOT NULL,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                               published_at TIMESTAMPTZ
);
CREATE INDEX outbox_events_unpublished_idx ON outbox_events(created_at) WHERE published_at IS NULL;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
DECLARE table_name TEXT;
BEGIN
    FOREACH table_name IN ARRAY ARRAY['users','identities','categories','products','cms','shopping_baskets','shopping_basket_items','reservations'] LOOP
        EXECUTE format('CREATE TRIGGER %I_updated_at BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION set_updated_at()', table_name, table_name);
END LOOP;
END $$;
-- +goose StatementEnd

-- +goose Down
DROP SCHEMA ecommerce CASCADE;
