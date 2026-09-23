package inventory

import "time"

// InventoryBalanceResponse represents a warehouse and variant balance response.
type InventoryBalanceResponse struct {
	WarehouseID string    `json:"warehouse_id"`
	VariantID   string    `json:"variant_id"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockMovementResponse represents an immutable stock movement response.
type StockMovementResponse struct {
	ID           string    `json:"id"`
	WarehouseID  string    `json:"warehouse_id"`
	VariantID    string    `json:"variant_id"`
	MovementType string    `json:"movement_type"`
	Quantity     int       `json:"quantity"`
	BalanceAfter int       `json:"balance_after"`
	Note         string    `json:"note"`
	CreatedAt    time.Time `json:"created_at"`
}

// AdjustmentResponse represents the balance and movement created by an adjustment.
type AdjustmentResponse struct {
	Balance  InventoryBalanceResponse `json:"balance"`
	Movement StockMovementResponse    `json:"movement"`
}
