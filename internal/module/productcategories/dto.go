package productcategories

import "time"

/*
======================================
Product Categories
======================================
*/

// ProductCategoryResponse represents a product category response.
type ProductCategoryResponse struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parent_id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}
