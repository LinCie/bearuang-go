package products

import (
	"context"
	"errors"
)

/*
======================================
Products
======================================
*/

var errInvalidProductStatus = errors.New(
	"status must be one of: draft, active, inactive, archived",
)

// Service contains product and product variant business operations.
type Service struct {
	repo Repository
}

// NewService creates a service for products and their variants backed by repo.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a product.
func (s *Service) Create(ctx context.Context, product *Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}

	return s.repo.Create(ctx, product)
}

// GetByID returns a product by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

// GetMany returns all non-deleted products.
func (s *Service) GetMany(ctx context.Context) ([]Product, error) {
	return s.repo.GetMany(ctx)
}

// Update updates a product.
func (s *Service) Update(ctx context.Context, product *Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}

	return s.repo.Update(ctx, product)
}

// Delete soft-deletes a product.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func validateProduct(product *Product) error {
	if product == nil {
		return nil
	}

	switch product.Status {
	case productStatusDraft, productStatusActive, productStatusInactive, productStatusArchived:
		return nil
	default:
		return errInvalidProductStatus
	}
}

/*
======================================
Variants
======================================
*/

var errInvalidVariantStatus = errors.New(
	"status must be one of: draft, active, inactive, archived",
)

// CreateVariant creates a product variant.
func (s *Service) CreateVariant(ctx context.Context, variant *ProductVariant) error {
	if err := validateVariant(variant); err != nil {
		return err
	}

	return s.repo.CreateVariant(ctx, variant)
}

// GetVariantByID returns a non-deleted product variant by ID for a product.
func (s *Service) GetVariantByID(ctx context.Context, productID, id string) (*ProductVariant, error) {
	return s.repo.GetVariantByID(ctx, productID, id)
}

// GetManyVariantsByProduct returns all non-deleted variants for a product.
func (s *Service) GetManyVariantsByProduct(ctx context.Context, productID string) ([]ProductVariant, error) {
	return s.repo.GetManyVariantsByProduct(ctx, productID)
}

// UpdateVariant updates a product variant.
func (s *Service) UpdateVariant(ctx context.Context, variant *ProductVariant) error {
	if err := validateVariant(variant); err != nil {
		return err
	}

	return s.repo.UpdateVariant(ctx, variant)
}

// DeleteVariant soft-deletes a non-deleted product variant by ID.
func (s *Service) DeleteVariant(ctx context.Context, productID, id string) error {
	return s.repo.DeleteVariant(ctx, productID, id)
}

func validateVariant(variant *ProductVariant) error {
	if variant == nil {
		return nil
	}

	switch variant.Status {
	case variantStatusDraft, variantStatusActive, variantStatusInactive, variantStatusArchived:
		return nil
	default:
		return errInvalidVariantStatus
	}
}
