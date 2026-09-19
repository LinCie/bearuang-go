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
	ID          string     `db:"id"`
	ParentID    *string    `db:"parent_id"`
	Name        string     `db:"name"`
	Slug        string     `db:"slug"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
