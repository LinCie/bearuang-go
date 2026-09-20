package inventory

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/database"
)

// Repository provides persistence for warehouses and inventory.
type Repository interface {
	CreateWarehouse(ctx context.Context, warehouse *Warehouse) error
	GetWarehouseByID(ctx context.Context, id string) (*Warehouse, error)
	GetManyWarehouses(ctx context.Context) ([]Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error

	GetWarehouseStocks(ctx context.Context, warehouseID string) ([]WarehouseStock, error)
	GetWarehouseStock(ctx context.Context, warehouseID, variantID string) (*WarehouseStock, error)
	AdjustStock(ctx context.Context, adjustment *StockAdjustment) error
	GetStockAdjustments(ctx context.Context, warehouseID string) ([]StockAdjustment, error)
}

// postgresRepository stores warehouses and inventory in PostgreSQL.
type postgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates an inventory repository backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *postgresRepository {
	return &postgresRepository{
		db: db,
	}
}

var _ Repository = (*postgresRepository)(nil)

/*
======================================
Warehouses
======================================
*/

// CreateWarehouse inserts a warehouse and populates its stored fields.
func (r *postgresRepository) CreateWarehouse(ctx context.Context, warehouse *Warehouse) error {
	if warehouse == nil {
		return errors.New("warehouse must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	warehouse.ID = database.GenerateID()

	query := `
		INSERT INTO WAREHOUSES (
			id,
			code,
			name,
			address,
			status
		) VALUES (:id, :code, :name, :address, :status)
		RETURNING
			id,
			code,
			name,
			address,
			status,
			created_at,
			updated_at,
			deleted_at`

	query, args, err := r.db.BindNamed(query, warehouse)
	if err != nil {
		return err
	}

	if err := r.db.GetContext(ctx, warehouse, query, args...); err != nil {
		return mapWarehouseConstraintError(err)
	}

	return nil
}

// GetWarehouseByID returns a non-deleted warehouse by ID.
func (r *postgresRepository) GetWarehouseByID(
	ctx context.Context,
	id string,
) (*Warehouse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			code,
			name,
			address,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM WAREHOUSES
		WHERE id = $1
			AND deleted_at IS NULL`

	var warehouse Warehouse
	if err := r.db.GetContext(ctx, &warehouse, query, id); err != nil {
		return nil, err
	}

	return &warehouse, nil
}

// GetManyWarehouses returns all non-deleted warehouses.
func (r *postgresRepository) GetManyWarehouses(ctx context.Context) ([]Warehouse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			code,
			name,
			address,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM WAREHOUSES
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC, id`

	warehouses := make([]Warehouse, 0)
	if err := r.db.SelectContext(ctx, &warehouses, query); err != nil {
		return nil, err
	}

	return warehouses, nil
}

// UpdateWarehouse updates a non-deleted warehouse and populates its stored fields.
func (r *postgresRepository) UpdateWarehouse(
	ctx context.Context,
	warehouse *Warehouse,
) error {
	if warehouse == nil {
		return errors.New("warehouse must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		UPDATE WAREHOUSES
		SET
			code = :code,
			name = :name,
			address = :address,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
			AND deleted_at IS NULL
		RETURNING
			id,
			code,
			name,
			address,
			status,
			created_at,
			updated_at,
			deleted_at`

	query, args, err := r.db.BindNamed(query, warehouse)
	if err != nil {
		return err
	}

	if err := r.db.GetContext(ctx, warehouse, query, args...); err != nil {
		return mapWarehouseConstraintError(err)
	}

	return nil
}

// DeleteWarehouse soft-deletes a warehouse only when it has no positive stock.
func (r *postgresRepository) DeleteWarehouse(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var warehouseID string
	if err := tx.GetContext(ctx, &warehouseID, `
		SELECT id
		FROM WAREHOUSES
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE`, id); err != nil {
		return err
	}

	var hasPositiveStock bool
	if err := tx.GetContext(ctx, &hasPositiveStock, `
		SELECT EXISTS (
			SELECT 1
			FROM WAREHOUSE_STOCKS
			WHERE warehouse_id = $1
				AND quantity > 0
		)`, id); err != nil {
		return err
	}
	if hasPositiveStock {
		return errWarehouseHasStock
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE WAREHOUSES
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

/*
======================================
Warehouse Stocks
======================================
*/

// GetWarehouseStocks returns all stock rows for a non-deleted warehouse.
func (r *postgresRepository) GetWarehouseStocks(
	ctx context.Context,
	warehouseID string,
) ([]WarehouseStock, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := r.ensureWarehouse(ctx, warehouseID); err != nil {
		return nil, err
	}

	query := `
		SELECT
			warehouse_id,
			variant_id,
			quantity,
			created_at,
			updated_at
		FROM WAREHOUSE_STOCKS
		WHERE warehouse_id = $1
		ORDER BY created_at DESC, variant_id`

	stocks := make([]WarehouseStock, 0)
	if err := r.db.SelectContext(ctx, &stocks, query, warehouseID); err != nil {
		return nil, err
	}

	return stocks, nil
}

// GetWarehouseStock returns one stock row for a non-deleted warehouse.
func (r *postgresRepository) GetWarehouseStock(
	ctx context.Context,
	warehouseID, variantID string,
) (*WarehouseStock, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := r.ensureWarehouse(ctx, warehouseID); err != nil {
		return nil, err
	}

	query := `
		SELECT
			warehouse_id,
			variant_id,
			quantity,
			created_at,
			updated_at
		FROM WAREHOUSE_STOCKS
		WHERE warehouse_id = $1
			AND variant_id = $2`

	var stock WarehouseStock
	if err := r.db.GetContext(ctx, &stock, query, warehouseID, variantID); err != nil {
		return nil, err
	}

	return &stock, nil
}

// AdjustStock atomically updates warehouse stock, aggregate variant stock,
// and append-only adjustment history in one transaction.
func (r *postgresRepository) AdjustStock(
	ctx context.Context,
	adjustment *StockAdjustment,
) error {
	if adjustment == nil {
		return errors.New("stock adjustment must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Every inventory adjustment locks its warehouse before its variant. This
	// serializes changes within a warehouse and gives all stock operations one
	// consistent lock ordering.
	var warehouseStatus string
	if err := tx.GetContext(ctx, &warehouseStatus, `
		SELECT status
		FROM WAREHOUSES
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE`, adjustment.WarehouseID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errWarehouseNotFound
		}
		return err
	}
	if warehouseStatus != warehouseStatusActive {
		return errWarehouseInactive
	}

	// Lock the aggregate row before changing either balance. The additive SQL
	// update below therefore cannot lose an update from another warehouse.
	var aggregateStock int
	if err := tx.GetContext(ctx, &aggregateStock, `
		SELECT stock
		FROM PRODUCT_VARIANTS
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE`, adjustment.VariantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errVariantNotFound
		}
		return err
	}

	var quantityBefore int
	stockRowExists := true
	if err := tx.GetContext(ctx, &quantityBefore, `
		SELECT quantity
		FROM WAREHOUSE_STOCKS
		WHERE warehouse_id = $1
			AND variant_id = $2
		FOR UPDATE`, adjustment.WarehouseID, adjustment.VariantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			stockRowExists = false
			quantityBefore = 0
		} else {
			return err
		}
	}

	var quantityAfter int
	if stockRowExists {
		if err := tx.GetContext(ctx, &quantityAfter, `
			UPDATE WAREHOUSE_STOCKS
			SET
				quantity = quantity + $3,
				updated_at = NOW()
			WHERE warehouse_id = $1
				AND variant_id = $2
				AND quantity + $3 >= 0
			RETURNING quantity`, adjustment.WarehouseID, adjustment.VariantID, adjustment.QuantityDelta); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errInsufficientWarehouseStock
			}
			return err
		}
	} else {
		if adjustment.QuantityDelta < 0 {
			return errInsufficientWarehouseStock
		}

		if err := tx.GetContext(ctx, &quantityAfter, `
			INSERT INTO WAREHOUSE_STOCKS (
				warehouse_id,
				variant_id,
				quantity
			) VALUES ($1, $2, $3)
			RETURNING quantity`, adjustment.WarehouseID, adjustment.VariantID, adjustment.QuantityDelta); err != nil {
			return err
		}
	}

	var aggregateAfter int
	if err := tx.GetContext(ctx, &aggregateAfter, `
		UPDATE PRODUCT_VARIANTS
		SET
			stock = stock + $1,
			updated_at = NOW()
		WHERE id = $2
			AND stock + $1 >= 0
		RETURNING stock`, adjustment.QuantityDelta, adjustment.VariantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errInsufficientAggregateStock
		}
		return err
	}

	adjustment.ID = database.GenerateID()
	adjustment.QuantityBefore = quantityBefore
	adjustment.QuantityAfter = quantityAfter

	query := `
		INSERT INTO STOCK_ADJUSTMENTS (
			id,
			warehouse_id,
			variant_id,
			quantity_delta,
			quantity_before,
			quantity_after,
			reason
		) VALUES (
			:id,
			:warehouse_id,
			:variant_id,
			:quantity_delta,
			:quantity_before,
			:quantity_after,
			:reason
		)
		RETURNING
			id,
			warehouse_id,
			variant_id,
			quantity_delta,
			quantity_before,
			quantity_after,
			reason,
			created_at`

	query, args, err := sqlx.Named(query, adjustment)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)
	if err := tx.GetContext(ctx, adjustment, query, args...); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetStockAdjustments returns adjustment history for a non-deleted warehouse.
func (r *postgresRepository) GetStockAdjustments(
	ctx context.Context,
	warehouseID string,
) ([]StockAdjustment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := r.ensureWarehouse(ctx, warehouseID); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			warehouse_id,
			variant_id,
			quantity_delta,
			quantity_before,
			quantity_after,
			reason,
			created_at
		FROM STOCK_ADJUSTMENTS
		WHERE warehouse_id = $1
		ORDER BY created_at DESC, id`

	adjustments := make([]StockAdjustment, 0)
	if err := r.db.SelectContext(ctx, &adjustments, query, warehouseID); err != nil {
		return nil, err
	}

	return adjustments, nil
}

func (r *postgresRepository) ensureWarehouse(ctx context.Context, id string) error {
	var exists bool
	if err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1
			FROM WAREHOUSES
			WHERE id = $1
				AND deleted_at IS NULL
		)`, id); err != nil {
		return err
	}
	if !exists {
		return errWarehouseNotFound
	}

	return nil
}

func mapWarehouseConstraintError(err error) error {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) &&
		pgError.Code == "23505" &&
		pgError.ConstraintName == "warehouse_code_key" {
		return errWarehouseCodeExists
	}

	return err
}
