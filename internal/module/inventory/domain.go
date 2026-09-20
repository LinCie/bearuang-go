package inventory

import "time"

/*
======================================
Warehouses
======================================
*/

const (
	warehouseStatusActive   = "active"
	warehouseStatusInactive = "inactive"
	warehouseStatusArchived = "archived"
)

// Warehouse represents a warehouse location.
type Warehouse struct {
	ID        string     `db:"id"`
	Code      string     `db:"code"`
	Name      string     `db:"name"`
	Address   string     `db:"address"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

/*
======================================
Warehouse Stocks
======================================
*/

// WarehouseStock represents the quantity of a product variant at a warehouse.
type WarehouseStock struct {
	WarehouseID string    `db:"warehouse_id"`
	VariantID   string    `db:"variant_id"`
	Quantity    int       `db:"quantity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

/*
======================================
Stock Adjustments
======================================
*/

// StockAdjustment records an append-only change to warehouse stock.
type StockAdjustment struct {
	ID             string    `db:"id"`
	WarehouseID    string    `db:"warehouse_id"`
	VariantID      string    `db:"variant_id"`
	QuantityDelta  int       `db:"quantity_delta"`
	QuantityBefore int       `db:"quantity_before"`
	QuantityAfter  int       `db:"quantity_after"`
	Reason         string    `db:"reason"`
	CreatedAt      time.Time `db:"created_at"`
}
