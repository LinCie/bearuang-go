package inventory

import (
	"context"
	"errors"
)

var (
	errInvalidWarehouseStatus = errors.New(
		"status must be one of: active, inactive, archived",
	)
	errWarehouseCodeExists        = errors.New("warehouse code already exists")
	errWarehouseHasStock          = errors.New("warehouse has positive stock")
	errWarehouseInactive          = errors.New("warehouse must be active for stock adjustments")
	errInvalidAdjustmentDelta     = errors.New("quantity_delta must not be zero")
	errInsufficientWarehouseStock = errors.New("warehouse does not contain enough stock")
	errInsufficientAggregateStock = errors.New("aggregate stock is insufficient")
	errWarehouseNotFound          = errors.New("warehouse not found")
	errVariantNotFound            = errors.New("product variant not found")
)

// Service provides warehouse and inventory operations to the handler.
type Service interface {
	CreateWarehouse(ctx context.Context, warehouse *Warehouse) error
	GetWarehouseByID(ctx context.Context, id string) (*Warehouse, error)
	GetManyWarehouses(ctx context.Context) ([]Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error

	GetWarehouseStocks(ctx context.Context, warehouseID string) ([]WarehouseStock, error)
	GetWarehouseStock(ctx context.Context, warehouseID, variantID string) (*WarehouseStock, error)
	AdjustStock(ctx context.Context, adjustment *StockAdjustment) error
	GetStockAdjustments(ctx context.Context, warehouseID string) ([]StockAdjustment, error)
}

// service contains warehouse and inventory business operations.
type service struct {
	repo Repository
}

// NewService creates an inventory service backed by repo.
func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

// CreateWarehouse creates a warehouse.
func (s *service) CreateWarehouse(ctx context.Context, warehouse *Warehouse) error {
	if err := validateWarehouse(warehouse); err != nil {
		return err
	}

	return s.repo.CreateWarehouse(ctx, warehouse)
}

// GetWarehouseByID returns a non-deleted warehouse by ID.
func (s *service) GetWarehouseByID(ctx context.Context, id string) (*Warehouse, error) {
	return s.repo.GetWarehouseByID(ctx, id)
}

// GetManyWarehouses returns all non-deleted warehouses.
func (s *service) GetManyWarehouses(ctx context.Context) ([]Warehouse, error) {
	return s.repo.GetManyWarehouses(ctx)
}

// UpdateWarehouse updates a non-deleted warehouse.
func (s *service) UpdateWarehouse(ctx context.Context, warehouse *Warehouse) error {
	if err := validateWarehouse(warehouse); err != nil {
		return err
	}

	return s.repo.UpdateWarehouse(ctx, warehouse)
}

// DeleteWarehouse soft-deletes a warehouse without positive stock.
func (s *service) DeleteWarehouse(ctx context.Context, id string) error {
	return s.repo.DeleteWarehouse(ctx, id)
}

// GetWarehouseStocks returns all stock rows for a non-deleted warehouse.
func (s *service) GetWarehouseStocks(ctx context.Context, warehouseID string) ([]WarehouseStock, error) {
	return s.repo.GetWarehouseStocks(ctx, warehouseID)
}

// GetWarehouseStock returns one stock row for a non-deleted warehouse.
func (s *service) GetWarehouseStock(
	ctx context.Context,
	warehouseID, variantID string,
) (*WarehouseStock, error) {
	return s.repo.GetWarehouseStock(ctx, warehouseID, variantID)
}

// AdjustStock atomically adjusts warehouse and aggregate variant stock.
func (s *service) AdjustStock(ctx context.Context, adjustment *StockAdjustment) error {
	if adjustment == nil {
		return errors.New("stock adjustment must not be nil")
	}
	if adjustment.QuantityDelta == 0 {
		return errInvalidAdjustmentDelta
	}

	return s.repo.AdjustStock(ctx, adjustment)
}

// GetStockAdjustments returns append-only adjustment history for a warehouse.
func (s *service) GetStockAdjustments(
	ctx context.Context,
	warehouseID string,
) ([]StockAdjustment, error) {
	return s.repo.GetStockAdjustments(ctx, warehouseID)
}

func validateWarehouse(warehouse *Warehouse) error {
	if warehouse == nil {
		return nil
	}

	switch warehouse.Status {
	case warehouseStatusActive, warehouseStatusInactive, warehouseStatusArchived:
		return nil
	default:
		return errInvalidWarehouseStatus
	}
}
