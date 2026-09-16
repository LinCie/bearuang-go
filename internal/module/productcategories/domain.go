package productcategories

import "time"

/*
======================================
Product Categories
======================================
*/

const (
	productCategoryStatusDraft    = "draft"
	productCategoryStatusActive   = "active"
	productCategoryStatusInactive = "inactive"
	productCategoryStatusArchived = "archived"
)

// ProductCategory represents a product category.
type ProductCategory struct {
	ID          string     `db:"id" json:"id"`
	ParentID    *string    `db:"parent_id" json:"parent_id"`
	Name        string     `db:"name" json:"name"`
	Slug        string     `db:"slug" json:"slug"`
	Description string     `db:"description" json:"description"`
	Status      string     `db:"status" json:"status"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`
}
