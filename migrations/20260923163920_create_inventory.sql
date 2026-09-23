-- +goose Up
-- Existing product_variants.stock values are not allocated because they have no warehouse attribution.
CREATE TABLE inventory_balances (
    warehouse_id VARCHAR(26) NOT NULL,
    variant_id VARCHAR(26) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT inventory_balance_pk PRIMARY KEY (warehouse_id, variant_id),
    CONSTRAINT inventory_balance_warehouse_fk FOREIGN KEY (warehouse_id) REFERENCES warehouses (id),
    CONSTRAINT inventory_balance_variant_fk FOREIGN KEY (variant_id) REFERENCES product_variants (id),
    CONSTRAINT inventory_balance_quantity_check CHECK (quantity >= 0)
);
CREATE INDEX inventory_balance_variant_idx ON inventory_balances (variant_id);

CREATE TABLE stock_movements (
    id VARCHAR(26) PRIMARY KEY,
    warehouse_id VARCHAR(26) NOT NULL,
    variant_id VARCHAR(26) NOT NULL,
    movement_type VARCHAR NOT NULL,
    quantity INTEGER NOT NULL,
    balance_after INTEGER NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT stock_movement_warehouse_fk FOREIGN KEY (warehouse_id) REFERENCES warehouses (id),
    CONSTRAINT stock_movement_variant_fk FOREIGN KEY (variant_id) REFERENCES product_variants (id),
    CONSTRAINT stock_movement_type_check CHECK (movement_type IN ('adjustment_in', 'adjustment_out')),
    CONSTRAINT stock_movement_quantity_check CHECK (quantity > 0),
    CONSTRAINT stock_movement_balance_after_check CHECK (balance_after >= 0)
);
CREATE INDEX stock_movement_warehouse_created_idx ON stock_movements (warehouse_id, created_at DESC);
CREATE INDEX stock_movement_variant_created_idx ON stock_movements (variant_id, created_at DESC);

-- +goose Down
DROP TABLE stock_movements;
DROP TABLE inventory_balances;
