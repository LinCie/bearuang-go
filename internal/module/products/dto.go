package products

import "time"

// ProductDTO is the JSON representation of a product.
type ProductDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func productDTO(product *Product) *ProductDTO {
	if product == nil {
		return nil
	}

	return &ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func productDTOs(products []Product) []ProductDTO {
	result := make([]ProductDTO, len(products))
	for i := range products {
		result[i] = *productDTO(&products[i])
	}

	return result
}
