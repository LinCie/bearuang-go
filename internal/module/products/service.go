package products

import (
	"context"
	"errors"
)

var errInvalidProductStatus = errors.New(
	"status must be one of: draft, active, inactive, archived",
)

// Service contains product business operations.
type Service struct {
	repo Repository
}

// NewService creates a product service backed by repo.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a product.
func (s *Service) Create(ctx context.Context, product *Product) (*Product, error) {
	if err := validateProduct(product); err != nil {
		return nil, err
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
func (s *Service) Update(ctx context.Context, product *Product) (*Product, error) {
	if err := validateProduct(product); err != nil {
		return nil, err
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
