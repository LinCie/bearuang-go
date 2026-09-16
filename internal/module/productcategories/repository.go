package productcategories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Repository provides persistence for product categories.
type Repository interface {
	Create(ctx context.Context, category *ProductCategory) error
	GetByID(ctx context.Context, id string) (*ProductCategory, error)
	GetMany(ctx context.Context) ([]ProductCategory, error)
	Update(ctx context.Context, category *ProductCategory) error
	Delete(ctx context.Context, id string) error
}

// PostgresRepository stores product categories in PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a repository for product categories backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

var _ Repository = (*PostgresRepository)(nil)

// Create inserts a product category and populates its stored fields.
func (r *PostgresRepository) Create(ctx context.Context, category *ProductCategory) error {
	if category == nil {
		return errors.New("product category must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		INSERT INTO PRODUCT_CATEGORIES (
			id,
			parent_id,
			name,
			slug,
			description,
			status
		) VALUES (:id, :parent_id, :name, :slug, :description, :status)
		RETURNING
			id,
			parent_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at`

	query, args, err := r.db.BindNamed(query, category)
	if err != nil {
		return err
	}

	if err := r.db.GetContext(ctx, category, query, args...); err != nil {
		return err
	}

	return nil
}

// GetByID returns a non-deleted product category by ID.
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			parent_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM PRODUCT_CATEGORIES
		WHERE id = $1
			AND deleted_at IS NULL`

	category := new(ProductCategory)
	if err := r.db.GetContext(ctx, category, query, id); err != nil {
		return nil, err
	}

	return category, nil
}

// GetMany returns all non-deleted product categories.
func (r *PostgresRepository) GetMany(ctx context.Context) ([]ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			parent_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM PRODUCT_CATEGORIES
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC, id`

	categories := make([]ProductCategory, 0)
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, err
	}

	return categories, nil
}

// Update updates a non-deleted product category and populates its stored fields.
func (r *PostgresRepository) Update(ctx context.Context, category *ProductCategory) error {
	if category == nil {
		return errors.New("product category must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		UPDATE PRODUCT_CATEGORIES
		SET
			parent_id = :parent_id,
			name = :name,
			slug = :slug,
			description = :description,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
			AND deleted_at IS NULL
		RETURNING
			id,
			parent_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at,
			deleted_at`

	query, args, err := r.db.BindNamed(query, category)
	if err != nil {
		return err
	}

	if err := r.db.GetContext(ctx, category, query, args...); err != nil {
		return err
	}

	return nil
}

// Delete soft-deletes a non-deleted product category.
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		UPDATE PRODUCT_CATEGORIES
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
