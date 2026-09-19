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

// toProductCategoryResponse converts a product category to its response DTO.
func toProductCategoryResponse(category ProductCategory) ProductCategoryResponse {
	return ProductCategoryResponse{
		ID:          category.ID,
		ParentID:    category.ParentID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Status:      category.Status,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
		DeletedAt:   category.DeletedAt,
	}
}

// toProductCategoryResponses converts product categories to response DTOs.
func toProductCategoryResponses(categories []ProductCategory) []ProductCategoryResponse {
	response := make([]ProductCategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, toProductCategoryResponse(category))
	}

	return response
}
