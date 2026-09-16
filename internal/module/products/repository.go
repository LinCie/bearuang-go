package products

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Repository provides persistence for products and their variants.
type Repository interface {
	/*
		======================================
		Products
		======================================
	*/

	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	GetMany(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error

	/*
		======================================
		Variants
		======================================
	*/

	CreateVariant(ctx context.Context, variant *ProductVariant) error
	GetVariantByID(ctx context.Context, productID, id string) (*ProductVariant, error)
	GetManyVariantsByProduct(ctx context.Context, productID string) ([]ProductVariant, error)
	UpdateVariant(ctx context.Context, variant *ProductVariant) error
	DeleteVariant(ctx context.Context, productID, id string) error
}

// PostgresRepository stores products and their variants in PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a repository for products and their variants backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

var _ Repository = (*PostgresRepository)(nil)

// Create inserts a product and populates its stored fields.
func (r *PostgresRepository) Create(ctx context.Context, product *Product) error {
	if product == nil {
		return errors.New("product must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
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
		return err
	}

	if err := r.db.GetContext(ctx, product, query, args...); err != nil {
		return err
	}

	return nil
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

// Update updates a non-deleted product and populates its stored fields.
func (r *PostgresRepository) Update(ctx context.Context, product *Product) error {
	if product == nil {
		return errors.New("product must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
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
		return err
	}

	if err := r.db.GetContext(ctx, product, query, args...); err != nil {
		return err
	}

	return nil
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

/*
======================================
Variants
======================================
*/

// CreateVariant inserts a product variant and populates its stored fields.
func (r *PostgresRepository) CreateVariant(ctx context.Context, variant *ProductVariant) error {
	if variant == nil {
		return errors.New("product variant must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		INSERT INTO PRODUCT_VARIANTS (
			id,
			product_id,
			sku,
			name,
			price,
			stock,
			unit,
			status
		) VALUES (:id, :product_id, :sku, :name, :price, :stock, :unit, :status)
		RETURNING
			id,
			product_id,
			sku,
			name,
			price,
			stock,
			unit,
			status,
			created_at,
			updated_at,
			deleted_at`
	query, args, err := r.db.BindNamed(query, variant)
	if err != nil {
		return err
	}

	if err := r.db.GetContext(ctx, variant, query, args...); err != nil {
		return err
	}

	return nil
}

// GetVariantByID returns a non-deleted product variant by ID for a product.
func (r *PostgresRepository) GetVariantByID(
	ctx context.Context,
	productID, id string,
) (*ProductVariant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			product_id,
			sku,
			name,
			price,
			stock,
			unit,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM PRODUCT_VARIANTS
		WHERE product_id = $1
			AND id = $2
			AND deleted_at IS NULL`

	variant := new(ProductVariant)
	if err := r.db.GetContext(ctx, variant, query, productID, id); err != nil {
		return nil, err
	}

	return variant, nil
}

// GetManyVariantsByProduct returns all non-deleted variants for a product.
func (r *PostgresRepository) GetManyVariantsByProduct(
	ctx context.Context,
	productID string,
) ([]ProductVariant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			product_id,
			sku,
			name,
			price,
			stock,
			unit,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM PRODUCT_VARIANTS
		WHERE product_id = $1
			AND deleted_at IS NULL
		ORDER BY created_at DESC, id`

	variants := make([]ProductVariant, 0)
	if err := r.db.SelectContext(ctx, &variants, query, productID); err != nil {
		return nil, err
	}

	return variants, nil
}

// UpdateVariant updates a non-deleted product variant and populates its stored fields.
func (r *PostgresRepository) UpdateVariant(
	ctx context.Context,
	variant *ProductVariant,
) error {
	if variant == nil {
		return errors.New("product variant must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		UPDATE PRODUCT_VARIANTS
		SET
			sku = :sku,
			name = :name,
			price = :price,
			stock = :stock,
			unit = :unit,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
			AND product_id = :product_id
			AND deleted_at IS NULL
		RETURNING
			id,
			product_id,
			sku,
			name,
			price,
			stock,
			unit,
			status,
			created_at,
			updated_at,
			deleted_at`

	query, args, err := r.db.BindNamed(query, variant)
	if err != nil {
		return err
	}

	if err := r.db.GetContext(ctx, variant, query, args...); err != nil {
		return err
	}

	return nil
}

// DeleteVariant soft-deletes a non-deleted product variant.
func (r *PostgresRepository) DeleteVariant(ctx context.Context, productID, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		UPDATE PRODUCT_VARIANTS
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE product_id = $1
			AND id = $2
			AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, productID, id)
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
