package products

import "context"

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
	return s.repo.Create(ctx, product)
}

// GetByID returns a product by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

// GetMany returns all active products.
func (s *Service) GetMany(ctx context.Context) ([]Product, error) {
	return s.repo.GetMany(ctx)
}

// Update updates a product.
func (s *Service) Update(ctx context.Context, product *Product) (*Product, error) {
	return s.repo.Update(ctx, product)
}

// Delete soft-deletes a product.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
