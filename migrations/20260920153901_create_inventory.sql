-- +goose Up
CREATE TABLE warehouses (
    id VARCHAR(26) PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL DEFAULT '',
    status VARCHAR NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT warehouse_code_key UNIQUE (code),
    CONSTRAINT warehouse_status_check CHECK (status IN ('active', 'inactive', 'archived'))
);

CREATE TABLE warehouse_stocks (
    warehouse_id VARCHAR(26) NOT NULL,
    variant_id VARCHAR(26) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT warehouse_stock_pkey PRIMARY KEY (warehouse_id, variant_id),
    CONSTRAINT warehouse_stock_warehouse_fk FOREIGN KEY (warehouse_id)
        REFERENCES warehouses (id) ON DELETE RESTRICT,
    CONSTRAINT warehouse_stock_variant_fk FOREIGN KEY (variant_id)
        REFERENCES product_variants (id) ON DELETE RESTRICT,
    CONSTRAINT warehouse_stock_quantity_check CHECK (quantity >= 0)
);

CREATE INDEX warehouse_stock_variant_idx ON warehouse_stocks (variant_id);

CREATE TABLE stock_adjustments (
    id VARCHAR(26) PRIMARY KEY,
    warehouse_id VARCHAR(26) NOT NULL,
    variant_id VARCHAR(26) NOT NULL,
    quantity_delta INTEGER NOT NULL,
    quantity_before INTEGER NOT NULL,
    quantity_after INTEGER NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT stock_adjustment_warehouse_fk FOREIGN KEY (warehouse_id)
        REFERENCES warehouses (id) ON DELETE RESTRICT,
    CONSTRAINT stock_adjustment_variant_fk FOREIGN KEY (variant_id)
        REFERENCES product_variants (id) ON DELETE RESTRICT,
    CONSTRAINT stock_adjustment_quantity_before_check CHECK (quantity_before >= 0),
    CONSTRAINT stock_adjustment_quantity_after_check CHECK (quantity_after >= 0),
    CONSTRAINT stock_adjustment_delta_check CHECK (quantity_delta <> 0)
);

CREATE INDEX stock_adjustment_warehouse_idx
    ON stock_adjustments (warehouse_id, created_at DESC, id);
CREATE INDEX stock_adjustment_variant_idx ON stock_adjustments (variant_id);

-- Existing product_variants.stock values are retained. There is no defined
-- default warehouse, so existing positive stock cannot be allocated safely.
-- Such stock remains unallocated until an explicit inventory policy exists.

-- +goose Down
DROP TABLE stock_adjustments;
DROP TABLE warehouse_stocks;
DROP TABLE warehouses;
