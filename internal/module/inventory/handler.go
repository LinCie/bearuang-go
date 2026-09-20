package inventory

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"bearuang-go/internal/httpx"
)

/*
======================================
Warehouses
======================================
*/

// Handler handles HTTP requests for warehouses and inventory.
type Handler struct {
	service Service
}

/*
======================================
Warehouse Types
======================================
*/

type warehouseCreateInput struct {
	Code    string `json:"code" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address"`
	Status  string `json:"status" validate:"required,oneof=active inactive archived"`
}

type warehouseIDInput struct {
	ID string `validate:"required"`
}

/*
======================================
Warehouse Handlers
======================================
*/

// NewHandler creates a handler for warehouses and inventory backed by service.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetManyWarehouses returns all non-deleted warehouses.
func (h *Handler) GetManyWarehouses(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.service.GetManyWarehouses(r.Context())
	if err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, toWarehouseResponses(warehouses))
}

// GetWarehouseByID returns a non-deleted warehouse by ID.
func (h *Handler) GetWarehouseByID(w http.ResponseWriter, r *http.Request) {
	id, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	warehouse, err := h.service.GetWarehouseByID(r.Context(), id)
	if err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, toWarehouseResponse(*warehouse))
}

// CreateWarehouse creates a warehouse.
func (h *Handler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var input warehouseCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	warehouse := Warehouse{
		Code:    input.Code,
		Name:    input.Name,
		Address: input.Address,
		Status:  input.Status,
	}
	if err := h.service.CreateWarehouse(r.Context(), &warehouse); err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusCreated, toWarehouseResponse(warehouse))
}

// UpdateWarehouse updates a non-deleted warehouse by ID.
func (h *Handler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	id, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	var input warehouseCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	warehouse := Warehouse{
		ID:      id,
		Code:    input.Code,
		Name:    input.Name,
		Address: input.Address,
		Status:  input.Status,
	}
	if err := h.service.UpdateWarehouse(r.Context(), &warehouse); err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, toWarehouseResponse(warehouse))
}

// DeleteWarehouse soft-deletes a warehouse with no positive stock.
func (h *Handler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	id, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	if err := h.service.DeleteWarehouse(r.Context(), id); err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, nil)
}

/*
======================================
Warehouse Stock Types
======================================
*/

type warehouseStockPathInput struct {
	WarehouseID string `json:"warehouse_id" validate:"required"`
	VariantID   string `json:"variant_id" validate:"required"`
}

/*
======================================
Warehouse Stock Handlers
======================================
*/

// GetWarehouseStocks returns all location-specific stock for a warehouse.
func (h *Handler) GetWarehouseStocks(w http.ResponseWriter, r *http.Request) {
	warehouseID, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	stocks, err := h.service.GetWarehouseStocks(r.Context(), warehouseID)
	if err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, toWarehouseStockResponses(stocks))
}

// GetWarehouseStock returns location-specific stock for one variant.
func (h *Handler) GetWarehouseStock(w http.ResponseWriter, r *http.Request) {
	input, ok := validateWarehouseStockPath(w, r)
	if !ok {
		return
	}

	stock, err := h.service.GetWarehouseStock(
		r.Context(),
		input.WarehouseID,
		input.VariantID,
	)
	if err != nil {
		respondInventoryError(w, err, "warehouse stock not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, toWarehouseStockResponse(*stock))
}

/*
======================================
Stock Adjustment Types
======================================
*/

type stockAdjustmentCreateInput struct {
	VariantID     string `json:"variant_id" validate:"required"`
	QuantityDelta int    `json:"quantity_delta"`
	Reason        string `json:"reason"`
}

/*
======================================
Stock Adjustment Handlers
======================================
*/

// CreateAdjustment atomically adjusts warehouse and aggregate variant stock.
func (h *Handler) CreateAdjustment(w http.ResponseWriter, r *http.Request) {
	warehouseID, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	var input stockAdjustmentCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	adjustment := StockAdjustment{
		WarehouseID:   warehouseID,
		VariantID:     input.VariantID,
		QuantityDelta: input.QuantityDelta,
		Reason:        input.Reason,
	}
	if err := h.service.AdjustStock(r.Context(), &adjustment); err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusCreated, toStockAdjustmentResponse(adjustment))
}

// GetStockAdjustments returns append-only adjustment history for a warehouse.
func (h *Handler) GetStockAdjustments(w http.ResponseWriter, r *http.Request) {
	warehouseID, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	adjustments, err := h.service.GetStockAdjustments(r.Context(), warehouseID)
	if err != nil {
		respondInventoryError(w, err, "warehouse not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, toStockAdjustmentResponses(adjustments))
}

/*
======================================
Helpers
======================================
*/

func validateWarehouseID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		id = strings.TrimSpace(chi.URLParam(r, "warehouse_id"))
	}

	if err := httpx.Validate(r.Context(), warehouseIDInput{ID: id}); err != nil {
		httpx.RespondInvalidBody(w, err)
		return "", false
	}

	return id, true
}

func validateWarehouseStockPath(
	w http.ResponseWriter,
	r *http.Request,
) (warehouseStockPathInput, bool) {
	input := warehouseStockPathInput{
		WarehouseID: strings.TrimSpace(chi.URLParam(r, "warehouse_id")),
		VariantID:   strings.TrimSpace(chi.URLParam(r, "variant_id")),
	}
	if err := httpx.Validate(r.Context(), input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return warehouseStockPathInput{}, false
	}

	return input, true
}

func respondInventoryError(w http.ResponseWriter, err error, notFoundMessage string) {
	switch {
	case errors.Is(err, errInvalidWarehouseStatus):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_status", err.Error())
	case errors.Is(err, errWarehouseCodeExists):
		httpx.RespondError(w, http.StatusConflict, "duplicate_code", err.Error())
	case errors.Is(err, errWarehouseHasStock):
		httpx.RespondError(w, http.StatusConflict, "warehouse_has_stock", err.Error())
	case errors.Is(err, errWarehouseInactive):
		httpx.RespondError(w, http.StatusConflict, "warehouse_inactive", err.Error())
	case errors.Is(err, errInvalidAdjustmentDelta):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_quantity_delta", err.Error())
	case errors.Is(err, errInsufficientWarehouseStock):
		httpx.RespondError(w, http.StatusConflict, "insufficient_stock", err.Error())
	case errors.Is(err, errInsufficientAggregateStock):
		httpx.RespondError(w, http.StatusConflict, "insufficient_aggregate_stock", err.Error())
	case errors.Is(err, errWarehouseNotFound):
		httpx.RespondError(w, http.StatusNotFound, "not_found", "warehouse not found")
	case errors.Is(err, errVariantNotFound):
		httpx.RespondError(w, http.StatusNotFound, "not_found", "product variant not found")
	case errors.Is(err, sql.ErrNoRows):
		httpx.RespondError(w, http.StatusNotFound, "not_found", notFoundMessage)
	default:
		httpx.RespondError(
			w,
			http.StatusInternalServerError,
			"internal_error",
			"internal server error",
		)
	}
}
