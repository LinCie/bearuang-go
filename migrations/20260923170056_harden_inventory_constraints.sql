-- +goose Up
-- Refuse to proceed when stock predates warehouse attribution or the stored
-- aggregate has otherwise drifted. Reconcile these variants explicitly before
-- applying this migration; never invent warehouse allocations or movements.
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM PRODUCT_VARIANTS AS variant
        LEFT JOIN (
            SELECT variant_id, SUM(quantity) AS quantity
            FROM INVENTORY_BALANCES
            GROUP BY variant_id
        ) AS balance ON balance.variant_id = variant.id
        WHERE variant.stock <> COALESCE(balance.quantity, 0)
    ) THEN
        RAISE EXCEPTION
            'cannot harden inventory constraints: product_variants.stock does not match warehouse balances; manually reconcile unattributed stock first';
    END IF;
END;
$$;
-- +goose StatementEnd

ALTER TABLE product_variants
ADD CONSTRAINT product_variant_stock_nonnegative_check
CHECK (stock >= 0);

-- +goose Down
ALTER TABLE product_variants
DROP CONSTRAINT product_variant_stock_nonnegative_check;
