package inventory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/database"
)

// Repository provides persistence for inventory balances and stock movements.
type Repository interface {
	GetBalance(ctx context.Context, warehouseID, variantID string) (*InventoryBalance, error)
	GetBalances(ctx context.Context, warehouseID, variantID *string) ([]InventoryBalance, error)
	Adjust(ctx context.Context, adjustment *Adjustment) (*InventoryBalance, *StockMovement, error)
	GetMovementByID(ctx context.Context, id string) (*StockMovement, error)
	GetMovements(ctx context.Context, warehouseID, variantID *string) ([]StockMovement, error)
}

// postgresRepository stores inventory in PostgreSQL.
type postgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates an inventory repository backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

var _ Repository = (*postgresRepository)(nil)

// GetBalance returns the current balance for a warehouse and variant.
func (r *postgresRepository) GetBalance(
	ctx context.Context,
	warehouseID, variantID string,
) (*InventoryBalance, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			warehouse_id,
			variant_id,
			quantity,
			created_at,
			updated_at
		FROM INVENTORY_BALANCES
		WHERE warehouse_id = $1
			AND variant_id = $2`

	var balance InventoryBalance
	if err := r.db.GetContext(ctx, &balance, query, warehouseID, variantID); err != nil {
		return nil, err
	}

	return &balance, nil
}

// GetBalances returns balances matching the optional warehouse and variant filters.
func (r *postgresRepository) GetBalances(
	ctx context.Context,
	warehouseID, variantID *string,
) ([]InventoryBalance, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			warehouse_id,
			variant_id,
			quantity,
			created_at,
			updated_at
		FROM INVENTORY_BALANCES`
	query, args := appendInventoryFilters(query, warehouseID, variantID)
	query += ` ORDER BY updated_at DESC, warehouse_id, variant_id`

	balances := make([]InventoryBalance, 0)
	if err := r.db.SelectContext(ctx, &balances, query, args...); err != nil {
		return nil, err
	}

	return balances, nil
}

// Adjust updates the balance, aggregate variant stock, and movement atomically.
func (r *postgresRepository) Adjust(
	ctx context.Context,
	adjustment *Adjustment,
) (*InventoryBalance, *StockMovement, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if adjustment == nil {
		return nil, nil, errors.New("adjustment must not be nil")
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var warehouseStatus string
	warehouseQuery := `
		SELECT status
		FROM WAREHOUSES
		WHERE id = $1
			AND deleted_at IS NULL
		FOR SHARE`
	if err := tx.GetContext(ctx, &warehouseStatus, warehouseQuery, adjustment.WarehouseID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, errWarehouseNotFound
		}
		return nil, nil, err
	}
	if warehouseStatus != "active" {
		return nil, nil, errWarehouseNotActive
	}

	variantExists, err := lockActiveProductAndVariant(ctx, tx, adjustment.VariantID)
	if err != nil {
		return nil, nil, err
	}
	if !variantExists {
		return nil, nil, errVariantNotFound
	}

	var balance InventoryBalance
	switch adjustment.Type {
	case movementTypeAdjustmentIn:
		query := `
			INSERT INTO INVENTORY_BALANCES (
				warehouse_id,
				variant_id,
				quantity
			) VALUES ($1, $2, $3)
			ON CONFLICT (warehouse_id, variant_id)
			DO UPDATE SET
				quantity = INVENTORY_BALANCES.quantity + EXCLUDED.quantity,
				updated_at = NOW()
			RETURNING
				warehouse_id,
				variant_id,
				quantity,
				created_at,
				updated_at`
		if err := tx.GetContext(
			ctx,
			&balance,
			query,
			adjustment.WarehouseID,
			adjustment.VariantID,
			adjustment.Quantity,
		); err != nil {
			return nil, nil, err
		}

	case movementTypeAdjustmentOut:
		query := `
			UPDATE INVENTORY_BALANCES
			SET
				quantity = quantity - $3,
				updated_at = NOW()
			WHERE warehouse_id = $1
				AND variant_id = $2
				AND quantity >= $3
			RETURNING
				warehouse_id,
				variant_id,
				quantity,
				created_at,
				updated_at`
		if err := tx.GetContext(
			ctx,
			&balance,
			query,
			adjustment.WarehouseID,
			adjustment.VariantID,
			adjustment.Quantity,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, errInsufficientStock
			}
			return nil, nil, err
		}

	default:
		return nil, nil, errInvalidAdjustmentType
	}

	if err := updateVariantStock(ctx, tx, adjustment); err != nil {
		return nil, nil, err
	}

	movement := StockMovement{
		ID:           database.GenerateID(),
		WarehouseID:  adjustment.WarehouseID,
		VariantID:    adjustment.VariantID,
		MovementType: adjustment.Type,
		Quantity:     adjustment.Quantity,
		BalanceAfter: balance.Quantity,
		Note:         adjustment.Note,
	}
	movementQuery := `
		INSERT INTO STOCK_MOVEMENTS (
			id,
			warehouse_id,
			variant_id,
			movement_type,
			quantity,
			balance_after,
			note
		) VALUES (
			:id,
			:warehouse_id,
			:variant_id,
			:movement_type,
			:quantity,
			:balance_after,
			:note
		)
		RETURNING
			id,
			warehouse_id,
			variant_id,
			movement_type,
			quantity,
			balance_after,
			note,
			created_at`
	movementQuery, args, err := tx.BindNamed(movementQuery, &movement)
	if err != nil {
		return nil, nil, err
	}
	if err := tx.GetContext(ctx, &movement, movementQuery, args...); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return &balance, &movement, nil
}

// GetMovementByID returns a movement by ID.
func (r *postgresRepository) GetMovementByID(ctx context.Context, id string) (*StockMovement, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			warehouse_id,
			variant_id,
			movement_type,
			quantity,
			balance_after,
			note,
			created_at
		FROM STOCK_MOVEMENTS
		WHERE id = $1`

	var movement StockMovement
	if err := r.db.GetContext(ctx, &movement, query, id); err != nil {
		return nil, err
	}

	return &movement, nil
}

// GetMovements returns movement history matching the optional filters, newest first.
func (r *postgresRepository) GetMovements(
	ctx context.Context,
	warehouseID, variantID *string,
) ([]StockMovement, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			warehouse_id,
			variant_id,
			movement_type,
			quantity,
			balance_after,
			note,
			created_at
		FROM STOCK_MOVEMENTS`
	query, args := appendInventoryFilters(query, warehouseID, variantID)
	query += ` ORDER BY created_at DESC, id DESC`

	movements := make([]StockMovement, 0)
	if err := r.db.SelectContext(ctx, &movements, query, args...); err != nil {
		return nil, err
	}

	return movements, nil
}

func appendInventoryFilters(query string, warehouseID, variantID *string) (string, []any) {
	filters := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if warehouseID != nil {
		args = append(args, *warehouseID)
		filters = append(filters, fmt.Sprintf("warehouse_id = $%d", len(args)))
	}
	if variantID != nil {
		args = append(args, *variantID)
		filters = append(filters, fmt.Sprintf("variant_id = $%d", len(args)))
	}
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	return query, args
}

// lockActiveProductAndVariant serializes adjustments with product and variant deletion before balance changes.
func lockActiveProductAndVariant(
	ctx context.Context,
	tx *sqlx.Tx,
	variantID string,
) (bool, error) {
	var productID string
	productQuery := `
		SELECT product.id
		FROM PRODUCTS AS product
		INNER JOIN PRODUCT_VARIANTS AS variant
			ON variant.product_id = product.id
		WHERE variant.id = $1
			AND variant.deleted_at IS NULL
			AND product.deleted_at IS NULL
		FOR SHARE OF product`
	if err := tx.GetContext(ctx, &productID, productQuery, variantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	var lockedVariantID string
	variantQuery := `
		SELECT id
		FROM PRODUCT_VARIANTS
		WHERE product_id = $1
			AND id = $2
			AND deleted_at IS NULL
		FOR NO KEY UPDATE`
	if err := tx.GetContext(ctx, &lockedVariantID, variantQuery, productID, variantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func activeVariantExists(ctx context.Context, tx *sqlx.Tx, variantID string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM PRODUCT_VARIANTS AS variant
			INNER JOIN PRODUCTS AS product
				ON product.id = variant.product_id
			WHERE variant.id = $1
				AND variant.deleted_at IS NULL
				AND product.deleted_at IS NULL
		)`
	if err := tx.GetContext(ctx, &exists, query, variantID); err != nil {
		return false, err
	}

	return exists, nil
}

func updateVariantStock(ctx context.Context, tx *sqlx.Tx, adjustment *Adjustment) error {
	query := `
		UPDATE PRODUCT_VARIANTS
		SET
			stock = stock + $2,
			updated_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL`
	args := []any{adjustment.VariantID, adjustment.Quantity}
	if adjustment.Type == movementTypeAdjustmentOut {
		query = `
			UPDATE PRODUCT_VARIANTS
			SET
				stock = stock - $2,
				updated_at = NOW()
			WHERE id = $1
				AND deleted_at IS NULL
				AND stock >= $2`
	}
	query += ` RETURNING stock`

	var stock int
	if err := tx.GetContext(ctx, &stock, query, args...); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		exists, existsErr := activeVariantExists(ctx, tx, adjustment.VariantID)
		if existsErr != nil {
			return existsErr
		}
		if !exists {
			return errVariantNotFound
		}
		return errInsufficientStock
	}

	return nil
}
