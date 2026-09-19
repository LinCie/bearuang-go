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

// toProductResponse converts a product to its response DTO.
func toProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		DeletedAt:   product.DeletedAt,
	}
}

// toProductResponses converts products to response DTOs.
func toProductResponses(products []Product) []ProductResponse {
	response := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, toProductResponse(product))
	}

	return response
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

// toProductVariantResponse converts a product variant to its response DTO.
func toProductVariantResponse(variant ProductVariant) ProductVariantResponse {
	return ProductVariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		SKU:       variant.SKU,
		Name:      variant.Name,
		Price:     variant.Price,
		Stock:     variant.Stock,
		Unit:      variant.Unit,
		Status:    variant.Status,
		CreatedAt: variant.CreatedAt,
		UpdatedAt: variant.UpdatedAt,
		DeletedAt: variant.DeletedAt,
	}
}

// toProductVariantResponses converts product variants to response DTOs.
func toProductVariantResponses(variants []ProductVariant) []ProductVariantResponse {
	response := make([]ProductVariantResponse, 0, len(variants))
	for _, variant := range variants {
		response = append(response, toProductVariantResponse(variant))
	}

	return response
}
