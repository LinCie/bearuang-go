package productcategories

import (
	"context"
	"errors"
)

// Repository provides persistence for product categories.
type Repository interface {
	Create(ctx context.Context, category *ProductCategory) error
	GetByID(ctx context.Context, id string) (*ProductCategory, error)
	GetMany(ctx context.Context) ([]ProductCategory, error)
	Update(ctx context.Context, category *ProductCategory) error
	Delete(ctx context.Context, id string) error
}

var errInvalidProductCategoryStatus = errors.New(
	"status must be one of: draft, active, inactive, archived",
)

// service contains product category business operations.
type service struct {
	repo Repository
}

// NewService creates a service for product categories backed by repo.
func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

// Create creates a product category.
func (s *service) Create(ctx context.Context, category *ProductCategory) error {
	if err := validateProductCategory(category); err != nil {
		return err
	}

	return s.repo.Create(ctx, category)
}

// GetByID returns a product category by ID.
func (s *service) GetByID(ctx context.Context, id string) (*ProductCategory, error) {
	return s.repo.GetByID(ctx, id)
}

// GetMany returns all non-deleted product categories.
func (s *service) GetMany(ctx context.Context) ([]ProductCategory, error) {
	return s.repo.GetMany(ctx)
}

// Update updates a product category.
func (s *service) Update(ctx context.Context, category *ProductCategory) error {
	if err := validateProductCategory(category); err != nil {
		return err
	}

	return s.repo.Update(ctx, category)
}

// Delete soft-deletes a product category.
func (s *service) Delete(ctx context.Context, id string) error {
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
