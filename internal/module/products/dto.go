package products

import "time"

/*
======================================
Products
======================================
*/

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

	dto := productDTOValue(*product)
	return &dto
}

func productDTOValue(product Product) ProductDTO {
	return ProductDTO{
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
		result[i] = productDTOValue(products[i])
	}

	return result
}

/*
======================================
Variants
======================================
*/

// ProductVariantDTO is the JSON representation of a product variant.
type ProductVariantDTO struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Unit      string    `json:"unit"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func productVariantDTO(variant *ProductVariant) *ProductVariantDTO {
	if variant == nil {
		return nil
	}

	dto := productVariantDTOValue(*variant)
	return &dto
}

func productVariantDTOValue(variant ProductVariant) ProductVariantDTO {
	return ProductVariantDTO{
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
	}
}

func productVariantDTOs(variants []ProductVariant) []ProductVariantDTO {
	result := make([]ProductVariantDTO, len(variants))
	for i := range variants {
		result[i] = productVariantDTOValue(variants[i])
	}

	return result
}
