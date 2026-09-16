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
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Slug        string     `db:"slug" json:"slug"`
	Description string     `db:"description" json:"description"`
	Status      string     `db:"status" json:"status"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`
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
	ID        string     `db:"id" json:"id"`
	ProductID string     `db:"product_id" json:"product_id"`
	SKU       string     `db:"sku" json:"sku"`
	Name      string     `db:"name" json:"name"`
	Price     float64    `db:"price" json:"price"`
	Stock     int        `db:"stock" json:"stock"`
	Unit      string     `db:"unit" json:"unit"`
	Status    string     `db:"status" json:"status"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
