package inventory

import "time"

/*
======================================
Warehouses
======================================
*/

// WarehouseResponse represents a warehouse response.
type WarehouseResponse struct {
	ID        string     `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// toWarehouseResponse converts a warehouse to its response DTO.
func toWarehouseResponse(warehouse Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID:        warehouse.ID,
		Code:      warehouse.Code,
		Name:      warehouse.Name,
		Address:   warehouse.Address,
		Status:    warehouse.Status,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
		DeletedAt: warehouse.DeletedAt,
	}
}

// toWarehouseResponses converts warehouses to response DTOs.
func toWarehouseResponses(warehouses []Warehouse) []WarehouseResponse {
	response := make([]WarehouseResponse, 0, len(warehouses))
	for _, warehouse := range warehouses {
		response = append(response, toWarehouseResponse(warehouse))
	}

	return response
}

/*
======================================
Warehouse Stocks
======================================
*/

// WarehouseStockResponse represents a warehouse stock response.
type WarehouseStockResponse struct {
	WarehouseID string    `json:"warehouse_id"`
	VariantID   string    `json:"variant_id"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// toWarehouseStockResponse converts warehouse stock to its response DTO.
func toWarehouseStockResponse(stock WarehouseStock) WarehouseStockResponse {
	return WarehouseStockResponse{
		WarehouseID: stock.WarehouseID,
		VariantID:   stock.VariantID,
		Quantity:    stock.Quantity,
		CreatedAt:   stock.CreatedAt,
		UpdatedAt:   stock.UpdatedAt,
	}
}

// toWarehouseStockResponses converts warehouse stocks to response DTOs.
func toWarehouseStockResponses(stocks []WarehouseStock) []WarehouseStockResponse {
	response := make([]WarehouseStockResponse, 0, len(stocks))
	for _, stock := range stocks {
		response = append(response, toWarehouseStockResponse(stock))
	}

	return response
}

/*
======================================
Stock Adjustments
======================================
*/

// StockAdjustmentResponse represents a stock adjustment response.
type StockAdjustmentResponse struct {
	ID             string    `json:"id"`
	WarehouseID    string    `json:"warehouse_id"`
	VariantID      string    `json:"variant_id"`
	QuantityDelta  int       `json:"quantity_delta"`
	QuantityBefore int       `json:"quantity_before"`
	QuantityAfter  int       `json:"quantity_after"`
	Reason         string    `json:"reason"`
	CreatedAt      time.Time `json:"created_at"`
}

// toStockAdjustmentResponse converts a stock adjustment to its response DTO.
func toStockAdjustmentResponse(adjustment StockAdjustment) StockAdjustmentResponse {
	return StockAdjustmentResponse{
		ID:             adjustment.ID,
		WarehouseID:    adjustment.WarehouseID,
		VariantID:      adjustment.VariantID,
		QuantityDelta:  adjustment.QuantityDelta,
		QuantityBefore: adjustment.QuantityBefore,
		QuantityAfter:  adjustment.QuantityAfter,
		Reason:         adjustment.Reason,
		CreatedAt:      adjustment.CreatedAt,
	}
}

// toStockAdjustmentResponses converts stock adjustments to response DTOs.
func toStockAdjustmentResponses(adjustments []StockAdjustment) []StockAdjustmentResponse {
	response := make([]StockAdjustmentResponse, 0, len(adjustments))
	for _, adjustment := range adjustments {
		response = append(response, toStockAdjustmentResponse(adjustment))
	}

	return response
}
