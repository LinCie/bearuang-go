package inventory

import (
	"context"
	"strings"
)

// Service provides inventory operations to the handler.
type Service interface {
	GetBalance(ctx context.Context, warehouseID, variantID string) (*InventoryBalance, error)
	GetBalances(ctx context.Context, warehouseID, variantID *string) ([]InventoryBalance, error)
	Adjust(ctx context.Context, adjustment *Adjustment) (*InventoryBalance, *StockMovement, error)
	GetMovementByID(ctx context.Context, id string) (*StockMovement, error)
	GetMovements(ctx context.Context, warehouseID, variantID *string) ([]StockMovement, error)
}

// service contains inventory business operations.
type service struct {
	repo Repository
}

// NewService creates an inventory service backed by repo.
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

// GetBalance returns the current balance for a warehouse and variant.
func (s *service) GetBalance(
	ctx context.Context,
	warehouseID, variantID string,
) (*InventoryBalance, error) {
	return s.repo.GetBalance(ctx, warehouseID, variantID)
}

// GetBalances returns balances matching the optional filters.
func (s *service) GetBalances(
	ctx context.Context,
	warehouseID, variantID *string,
) ([]InventoryBalance, error) {
	return s.repo.GetBalances(ctx, warehouseID, variantID)
}

// Adjust applies a validated stock adjustment.
func (s *service) Adjust(
	ctx context.Context,
	adjustment *Adjustment,
) (*InventoryBalance, *StockMovement, error) {
	if adjustment == nil {
		return nil, nil, errInvalidAdjustmentQuantity
	}

	adjustment.WarehouseID = strings.TrimSpace(adjustment.WarehouseID)
	adjustment.VariantID = strings.TrimSpace(adjustment.VariantID)
	if adjustment.WarehouseID == "" {
		return nil, nil, errInvalidWarehouseID
	}
	if adjustment.VariantID == "" {
		return nil, nil, errInvalidVariantID
	}
	if adjustment.Type != movementTypeAdjustmentIn && adjustment.Type != movementTypeAdjustmentOut {
		return nil, nil, errInvalidAdjustmentType
	}
	if adjustment.Quantity <= 0 {
		return nil, nil, errInvalidAdjustmentQuantity
	}

	return s.repo.Adjust(ctx, adjustment)
}

// GetMovementByID returns a movement by ID.
func (s *service) GetMovementByID(ctx context.Context, id string) (*StockMovement, error) {
	return s.repo.GetMovementByID(ctx, id)
}

// GetMovements returns movement history matching the optional filters.
func (s *service) GetMovements(
	ctx context.Context,
	warehouseID, variantID *string,
) ([]StockMovement, error) {
	return s.repo.GetMovements(ctx, warehouseID, variantID)
}
