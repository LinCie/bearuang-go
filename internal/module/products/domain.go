package products

import "time"

/*
======================================
Products
======================================
*/

const (
	productStatusDraft    = "draft"
	productStatusActive   = "active"
	productStatusInactive = "inactive"
	productStatusArchived = "archived"
)

// Product represents a product.
type Product struct {
	ID          string     `db:"id"`
	CategoryID  *string    `db:"category_id"`
	Name        string     `db:"name"`
	Slug        string     `db:"slug"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

/*
======================================
Variants
======================================
*/

const (
	variantStatusDraft    = "draft"
	variantStatusActive   = "active"
	variantStatusInactive = "inactive"
	variantStatusArchived = "archived"
)

// ProductVariant represents a product variant.
type ProductVariant struct {
	ID        string  `db:"id"`
	ProductID string  `db:"product_id"`
	SKU       string  `db:"sku"`
	Name      string  `db:"name"`
	Price     float64 `db:"price"`
	// Stock is the denormalized aggregate maintained by inventory operations.
	Stock     int        `db:"stock"`
	Unit      string     `db:"unit"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
