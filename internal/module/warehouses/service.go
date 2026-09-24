package warehouses

import (
	"context"
	"errors"
	"strings"
)

var (
	errInvalidWarehouseCode   = errors.New("code must not be empty")
	errInvalidWarehouseName   = errors.New("name must not be empty")
	errInvalidWarehouseStatus = errors.New(
		"status must be one of: active, inactive, archived",
	)
	errDuplicateWarehouseCode = errors.New("warehouse code already exists")
	errWarehouseInUse         = errors.New("warehouse has inventory history and must be archived instead")
)

// service contains warehouse business operations.
type service struct {
	repo Repository
}

// NewService creates a service for warehouses backed by repo.
func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

// Create creates a warehouse.
func (s *service) Create(ctx context.Context, warehouse *Warehouse) error {
	if err := normalizeAndValidateWarehouse(warehouse); err != nil {
		return err
	}

	return s.repo.Create(ctx, warehouse)
}

// GetByID returns a warehouse by ID.
func (s *service) GetByID(ctx context.Context, id string) (*Warehouse, error) {
	return s.repo.GetByID(ctx, id)
}

// GetMany returns all non-deleted warehouses.
func (s *service) GetMany(ctx context.Context) ([]Warehouse, error) {
	return s.repo.GetMany(ctx)
}

// Update updates a warehouse.
func (s *service) Update(ctx context.Context, warehouse *Warehouse) error {
	if err := normalizeAndValidateWarehouse(warehouse); err != nil {
		return err
	}

	return s.repo.Update(ctx, warehouse)
}

// Delete soft-deletes a warehouse.
func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func normalizeAndValidateWarehouse(warehouse *Warehouse) error {
	if warehouse == nil {
		return nil
	}

	warehouse.Code = strings.TrimSpace(warehouse.Code)
	warehouse.Name = strings.TrimSpace(warehouse.Name)
	warehouse.Description = strings.TrimSpace(warehouse.Description)
	warehouse.Status = strings.TrimSpace(warehouse.Status)
	if warehouse.Address != nil {
		address := strings.TrimSpace(*warehouse.Address)
		warehouse.Address = &address
	}

	if warehouse.Code == "" {
		return errInvalidWarehouseCode
	}
	if warehouse.Name == "" {
		return errInvalidWarehouseName
	}

	switch warehouse.Status {
	case warehouseStatusActive, warehouseStatusInactive, warehouseStatusArchived:
		return nil
	default:
		return errInvalidWarehouseStatus
	}
}
