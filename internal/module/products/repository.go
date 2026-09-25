package products

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/database"
)

// postgresRepository stores products and their variants in PostgreSQL.
type postgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a repository for products and their variants backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *postgresRepository {
	return &postgresRepository{
		db: db,
	}
}

var _ Repository = (*postgresRepository)(nil)

// Create inserts a product and populates its stored fields.
func (r *postgresRepository) Create(ctx context.Context, product *Product) error {
	if product == nil {
		return errors.New("product must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	product.ID = database.GenerateID()

	query := `
		INSERT INTO PRODUCTS (
			id,
			category_id,
			name,
			slug,
			description,
			status
		) VALUES (:id, :category_id, :name, :slug, :description, :status)
		RETURNING
			id,
			category_id,
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
func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			category_id,
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

	var product Product
	if err := r.db.GetContext(ctx, &product, query, id); err != nil {
		return nil, err
	}

	return &product, nil
}

// GetMany returns all non-deleted products.
func (r *postgresRepository) GetMany(ctx context.Context) ([]Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			category_id,
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
func (r *postgresRepository) Update(ctx context.Context, product *Product) error {
	if product == nil {
		return errors.New("product must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		UPDATE PRODUCTS
		SET
			category_id = :category_id,
			name = :name,
			slug = :slug,
			description = :description,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
			AND deleted_at IS NULL
		RETURNING
			id,
			category_id,
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

// Delete soft-deletes a product with no inventory-relevant variants.
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

	var productID string
	lockQuery := `
		SELECT id
		FROM PRODUCTS
		WHERE id = $1
			AND deleted_at IS NULL
		FOR UPDATE`
	if err := tx.GetContext(ctx, &productID, lockQuery, id); err != nil {
		return err
	}

	var used bool
	usageQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM PRODUCT_VARIANTS variant
			WHERE variant.product_id = $1
				AND (
					variant.stock <> 0
					OR EXISTS (
						SELECT 1
						FROM INVENTORY_BALANCES balance
						WHERE balance.variant_id = variant.id
					)
					OR EXISTS (
						SELECT 1
						FROM STOCK_MOVEMENTS movement
						WHERE movement.variant_id = variant.id
					)
				)
		)`
	if err := tx.GetContext(ctx, &used, usageQuery, id); err != nil {
		return err
	}
	if used {
		return errProductInUse
	}

	deleteQuery := `
		UPDATE PRODUCTS
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

/*
======================================
Variants
======================================
*/

// CreateVariant inserts a product variant and populates its stored fields.
func (r *postgresRepository) CreateVariant(ctx context.Context, variant *ProductVariant) error {
	if variant == nil {
		return errors.New("product variant must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	variant.ID = database.GenerateID()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var productID string
	lockQuery := `
		SELECT id
		FROM PRODUCTS
		WHERE id = $1
			AND deleted_at IS NULL
		FOR SHARE`
	if err := tx.GetContext(ctx, &productID, lockQuery, variant.ProductID); err != nil {
		return err
	}

	query := `
		INSERT INTO PRODUCT_VARIANTS (
			id,
			product_id,
			sku,
			name,
			price,
			unit,
			status
		) VALUES (:id, :product_id, :sku, :name, :price, :unit, :status)
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

	query, args, err := tx.BindNamed(query, variant)
	if err != nil {
		return err
	}
	if err := tx.GetContext(ctx, variant, query, args...); err != nil {
		return err
	}

	return tx.Commit()
}

// GetVariantByID returns a non-deleted product variant by ID for a product.
func (r *postgresRepository) GetVariantByID(
	ctx context.Context,
	productID, id string,
) (*ProductVariant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			variant.id,
			variant.product_id,
			variant.sku,
			variant.name,
			variant.price,
			variant.stock,
			variant.unit,
			variant.status,
			variant.created_at,
			variant.updated_at,
			variant.deleted_at
		FROM PRODUCT_VARIANTS AS variant
		INNER JOIN PRODUCTS AS product
			ON product.id = variant.product_id
		WHERE product.id = $1
			AND variant.id = $2
			AND product.deleted_at IS NULL
			AND variant.deleted_at IS NULL`

	var variant ProductVariant
	if err := r.db.GetContext(ctx, &variant, query, productID, id); err != nil {
		return nil, err
	}

	return &variant, nil
}

// GetManyVariantsByProduct returns all non-deleted variants for a product.
func (r *postgresRepository) GetManyVariantsByProduct(
	ctx context.Context,
	productID string,
) ([]ProductVariant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var lockedProductID string
	lockQuery := `
		SELECT id
		FROM PRODUCTS
		WHERE id = $1
			AND deleted_at IS NULL
		FOR SHARE`
	if err := tx.GetContext(ctx, &lockedProductID, lockQuery, productID); err != nil {
		return nil, err
	}

	query := `
		SELECT
			variant.id,
			variant.product_id,
			variant.sku,
			variant.name,
			variant.price,
			variant.stock,
			variant.unit,
			variant.status,
			variant.created_at,
			variant.updated_at,
			variant.deleted_at
		FROM PRODUCT_VARIANTS AS variant
		INNER JOIN PRODUCTS AS product
			ON product.id = variant.product_id
		WHERE product.id = $1
			AND product.deleted_at IS NULL
			AND variant.deleted_at IS NULL
		ORDER BY variant.created_at DESC, variant.id`

	variants := make([]ProductVariant, 0)
	if err := tx.SelectContext(ctx, &variants, query, productID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return variants, nil
}

// UpdateVariant updates a non-deleted product variant and populates its stored fields.
func (r *postgresRepository) UpdateVariant(
	ctx context.Context,
	variant *ProductVariant,
) error {
	if variant == nil {
		return errors.New("product variant must not be nil")
	}
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

	var productID string
	lockQuery := `
		SELECT id
		FROM PRODUCTS
		WHERE id = $1
			AND deleted_at IS NULL
		FOR SHARE`
	if err := tx.GetContext(ctx, &productID, lockQuery, variant.ProductID); err != nil {
		return err
	}

	query := `
		UPDATE PRODUCT_VARIANTS
		SET
			sku = :sku,
			name = :name,
			price = :price,
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

	query, args, err := tx.BindNamed(query, variant)
	if err != nil {
		return err
	}
	if err := tx.GetContext(ctx, variant, query, args...); err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteVariant soft-deletes a variant with no inventory state or history.
func (r *postgresRepository) DeleteVariant(ctx context.Context, productID, id string) error {
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

	var lockedProductID string
	productLockQuery := `
		SELECT id
		FROM PRODUCTS
		WHERE id = $1
			AND deleted_at IS NULL
		FOR SHARE`
	if err := tx.GetContext(ctx, &lockedProductID, productLockQuery, productID); err != nil {
		return err
	}

	var stock int
	lockQuery := `
		SELECT stock
		FROM PRODUCT_VARIANTS
		WHERE product_id = $1
			AND id = $2
			AND deleted_at IS NULL
		FOR UPDATE`
	if err := tx.GetContext(ctx, &stock, lockQuery, productID, id); err != nil {
		return err
	}

	var used bool
	usageQuery := `
		SELECT
			$1 <> 0
			OR EXISTS (
				SELECT 1
				FROM INVENTORY_BALANCES
				WHERE variant_id = $2
			) OR EXISTS (
				SELECT 1
				FROM STOCK_MOVEMENTS
				WHERE variant_id = $2
			)`
	if err := tx.GetContext(ctx, &used, usageQuery, stock, id); err != nil {
		return err
	}
	if used {
		return errVariantInUse
	}

	deleteQuery := `
		UPDATE PRODUCT_VARIANTS
		SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE product_id = $1
			AND id = $2
			AND deleted_at IS NULL`
	result, err := tx.ExecContext(ctx, deleteQuery, productID, id)
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
