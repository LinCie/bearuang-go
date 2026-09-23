package products

import "time"

/*
======================================
Products
======================================
*/

// ProductResponse represents a product response.
type ProductResponse struct {
	ID          string     `json:"id"`
	CategoryID  *string    `json:"category_id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

/*
======================================
Variants
======================================
*/

// ProductVariantResponse represents a product variant response.
type ProductVariantResponse struct {
	ID        string     `json:"id"`
	ProductID string     `json:"product_id"`
	SKU       string     `json:"sku"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	Stock     int        `json:"stock"`
	Unit      string     `json:"unit"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

