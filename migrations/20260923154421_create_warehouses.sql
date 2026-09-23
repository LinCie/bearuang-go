-- +goose Up
CREATE TABLE warehouses (
    id VARCHAR(26) PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    address TEXT,
    status VARCHAR NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT warehouse_code_key UNIQUE (code)
);

-- +goose Down
DROP TABLE warehouses;
