package products

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Repository provides persistence for products.
type Repository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	GetMany(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	Delete(ctx context.Context, id string) error
}

// PostgresRepository stores products in PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a product repository backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

var _ Repository = (*PostgresRepository)(nil)

// Create inserts a product and returns the stored row.
func (r *PostgresRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	if product == nil {
		return nil, errors.New("product must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO PRODUCTS (
			id,
			name,
			slug,
			description,
			status
		) VALUES (:id, :name, :slug, :description, :status)
		RETURNING
			id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at`
	query, args, err := r.db.BindNamed(query, product)
	if err != nil {
		return nil, err
	}

	stored := new(Product)
	if err := r.db.GetContext(ctx, stored, query, args...); err != nil {
		return nil, err
	}

	return stored, nil
}

// GetByID returns a non-deleted product by ID.
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM PRODUCTS
		WHERE id = $1
			AND deleted_at IS NULL`

	product := new(Product)
	if err := r.db.GetContext(ctx, product, query, id); err != nil {
		return nil, err
	}

	return product, nil
}

// GetMany returns all non-deleted products.
func (r *PostgresRepository) GetMany(ctx context.Context) ([]Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM PRODUCTS
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC, id`

	products := make([]Product, 0)
	if err := r.db.SelectContext(ctx, &products, query); err != nil {
		return nil, err
	}

	return products, nil
}

// Update updates a non-deleted product and returns the stored row.
func (r *PostgresRepository) Update(ctx context.Context, product *Product) (*Product, error) {
	if product == nil {
		return nil, errors.New("product must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		UPDATE PRODUCTS
		SET
			name = :name,
			slug = :slug,
			description = :description,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
			AND deleted_at IS NULL
		RETURNING
			id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at`

	query, args, err := r.db.BindNamed(query, product)
	if err != nil {
		return nil, err
	}

	stored := new(Product)
	if err := r.db.GetContext(ctx, stored, query, args...); err != nil {
		return nil, err
	}

	return stored, nil
}

// Delete soft-deletes a non-deleted product.
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		UPDATE PRODUCTS
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
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

	return nil
}
