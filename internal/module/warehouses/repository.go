package warehouses

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/database"
)

// Repository provides persistence for warehouses.
type Repository interface {
	Create(ctx context.Context, warehouse *Warehouse) error
	GetByID(ctx context.Context, id string) (*Warehouse, error)
	GetMany(ctx context.Context) ([]Warehouse, error)
	Update(ctx context.Context, warehouse *Warehouse) error
	Delete(ctx context.Context, id string) error
}

// postgresRepository stores warehouses in PostgreSQL.
type postgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a repository for warehouses backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *postgresRepository {
	return &postgresRepository{
		db: db,
	}
}

var _ Repository = (*postgresRepository)(nil)

// Create inserts a warehouse and populates its stored fields.
func (r *postgresRepository) Create(ctx context.Context, warehouse *Warehouse) error {
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
			description,
			address,
			status
		) VALUES (:id, :code, :name, :description, :address, :status)
		RETURNING
			id,
			code,
			name,
			description,
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
		return warehouseWriteError(err)
	}

	return nil
}

// GetByID returns a non-deleted warehouse by ID.
func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Warehouse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			code,
			name,
			description,
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

// GetMany returns all non-deleted warehouses.
func (r *postgresRepository) GetMany(ctx context.Context) ([]Warehouse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			code,
			name,
			description,
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

// Update updates a non-deleted warehouse and populates its stored fields.
func (r *postgresRepository) Update(ctx context.Context, warehouse *Warehouse) error {
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
			description = :description,
			address = :address,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
			AND deleted_at IS NULL
		RETURNING
			id,
			code,
			name,
			description,
			address,
			status,
			created_at,
			updated_at,
			deleted_at`

	if warehouse.Status == warehouseStatusActive {
		query, args, err := r.db.BindNamed(query, warehouse)
		if err != nil {
			return err
		}

		if err := r.db.GetContext(ctx, warehouse, query, args...); err != nil {
			return warehouseWriteError(err)
		}
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var warehouseID string
	lockQuery := `
		SELECT id
		FROM WAREHOUSES
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE`
	if err := tx.GetContext(ctx, &warehouseID, lockQuery, warehouse.ID); err != nil {
		return err
	}

	var hasStock bool
	stockQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM INVENTORY_BALANCES
			WHERE warehouse_id = $1
				AND quantity > 0
		)`
	if err := tx.GetContext(ctx, &hasStock, stockQuery, warehouse.ID); err != nil {
		return err
	}
	if hasStock {
		return errWarehouseHasStock
	}

	query, args, err := tx.BindNamed(query, warehouse)
	if err != nil {
		return err
	}
	if err := tx.GetContext(ctx, warehouse, query, args...); err != nil {
		return warehouseWriteError(err)
	}

	return tx.Commit()
}

// Delete soft-deletes an unused warehouse.
func (r *postgresRepository) Delete(ctx context.Context, id string) error {
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
	lockQuery := `
		SELECT id
		FROM WAREHOUSES
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE`
	if err := tx.GetContext(ctx, &warehouseID, lockQuery, id); err != nil {
		return err
	}

	var used bool
	usageQuery := `
		SELECT
			EXISTS (
				SELECT 1
				FROM INVENTORY_BALANCES
				WHERE warehouse_id = $1
			) OR EXISTS (
				SELECT 1
				FROM STOCK_MOVEMENTS
				WHERE warehouse_id = $1
			)`
	if err := tx.GetContext(ctx, &used, usageQuery, id); err != nil {
		return err
	}
	if used {
		return errWarehouseInUse
	}

	deleteQuery := `
		UPDATE WAREHOUSES
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL`
	result, err := tx.ExecContext(ctx, deleteQuery, id)
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

	return tx.Commit()
}

func warehouseWriteError(err error) error {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) &&
		pgError.Code == "23505" &&
		pgError.ConstraintName == "warehouse_code_key" {
		return errDuplicateWarehouseCode
	}

	return err
}
