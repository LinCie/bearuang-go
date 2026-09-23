package inventory

import (
	"errors"
	"time"
)

/*
======================================
Inventory Balances
======================================
*/

const (
	movementTypeAdjustmentIn  = "adjustment_in"
	movementTypeAdjustmentOut = "adjustment_out"
)

var (
	errWarehouseNotFound         = errors.New("warehouse not found")
	errWarehouseNotActive        = errors.New("warehouse is not active")
	errVariantNotFound           = errors.New("product variant not found")
	errInsufficientStock         = errors.New("insufficient stock")
	errInvalidAdjustmentType     = errors.New("unsupported adjustment type")
	errInvalidAdjustmentQuantity = errors.New("quantity must be greater than zero")
	errInvalidWarehouseID        = errors.New("warehouse_id must not be empty")
	errInvalidVariantID          = errors.New("variant_id must not be empty")
)

// InventoryBalance represents current stock for a warehouse and product variant.
type InventoryBalance struct {
	WarehouseID string    `db:"warehouse_id"`
	VariantID   string    `db:"variant_id"`
	Quantity    int       `db:"quantity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// StockMovement represents an immutable inventory change event.
type StockMovement struct {
	ID           string    `db:"id"`
	WarehouseID  string    `db:"warehouse_id"`
	VariantID    string    `db:"variant_id"`
	MovementType string    `db:"movement_type"`
	Quantity     int       `db:"quantity"`
	BalanceAfter int       `db:"balance_after"`
	Note         string    `db:"note"`
	CreatedAt    time.Time `db:"created_at"`
}

// Adjustment describes a manual increase or decrease in inventory.
type Adjustment struct {
	WarehouseID string
	VariantID   string
	Type        string
	Quantity    int
	Note        string
}
