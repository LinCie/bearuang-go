-- +goose Up
CREATE TABLE product_category (
    id VARCHAR(26) PRIMARY KEY,
    parent_id VARCHAR(26),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT product_category_parent_fk FOREIGN KEY (parent_id) REFERENCES product_category (id) ON DELETE
    SET NULL,
        CONSTRAINT product_category_slug_key UNIQUE (slug)
);
CREATE INDEX product_category_parent_idx ON product_category (parent_id);
CREATE TABLE product (
    id VARCHAR(26) PRIMARY KEY,
    category_id VARCHAR(26),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT product_category_fk FOREIGN KEY (category_id) REFERENCES product_category (id) ON DELETE
    SET NULL,
        CONSTRAINT product_slug_key UNIQUE (slug)
);
CREATE INDEX product_category_idx ON product (category_id);
CREATE TABLE product_variant (
    id VARCHAR(26) PRIMARY KEY,
    product_id VARCHAR(26) NOT NULL,
    sku TEXT NOT NULL,
    name TEXT NOT NULL,
    price NUMERIC(12, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    unit TEXT NOT NULL DEFAULT 'pcs',
    status VARCHAR NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT product_variant_product_fk FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE,
    CONSTRAINT product_variant_sku_key UNIQUE (sku)
);
CREATE INDEX product_variant_product_idx ON product_variant (product_id);
-- +goose Down
DROP TABLE product_variant;
DROP TABLE product;
DROP TABLE product_category;