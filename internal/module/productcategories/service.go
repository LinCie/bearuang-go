package productcategories

import (
	"context"
	"errors"
)

var errInvalidProductCategoryStatus = errors.New(
	"status must be one of: draft, active, inactive, archived",
)

// Service contains product category business operations.
type Service struct {
	repo Repository
}

// NewService creates a service for product categories backed by repo.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a product category.
func (s *Service) Create(ctx context.Context, category *ProductCategory) error {
	if err := validateProductCategory(category); err != nil {
		return err
	}

	return s.repo.Create(ctx, category)
}

// GetByID returns a product category by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*ProductCategory, error) {
	return s.repo.GetByID(ctx, id)
}

// GetMany returns all non-deleted product categories.
func (s *Service) GetMany(ctx context.Context) ([]ProductCategory, error) {
	return s.repo.GetMany(ctx)
}

// Update updates a product category.
func (s *Service) Update(ctx context.Context, category *ProductCategory) error {
	if err := validateProductCategory(category); err != nil {
		return err
	}

	return s.repo.Update(ctx, category)
}

// Delete soft-deletes a product category.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func validateProductCategory(category *ProductCategory) error {
	if category == nil {
		return nil
	}

	switch category.Status {
	case productCategoryStatusDraft,
		productCategoryStatusActive,
		productCategoryStatusInactive,
		productCategoryStatusArchived:
		return nil
	default:
		return errInvalidProductCategoryStatus
	}
}
