package inventory

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"bearuang-go/internal/httpx"
)

// Handler handles HTTP requests for inventory balances and movements.
type Handler struct {
	service Service
}

type adjustmentInput struct {
	WarehouseID string `json:"warehouse_id" validate:"required"`
	VariantID   string `json:"variant_id" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=adjustment_in adjustment_out"`
	Quantity    int    `json:"quantity" validate:"required,gt=0"`
	Note        string `json:"note"`
}

type inventoryBalancePathInput struct {
	WarehouseID string `json:"warehouse_id" validate:"required"`
	VariantID   string `json:"variant_id" validate:"required"`
}

type movementIDInput struct {
	ID string `json:"id" validate:"required"`
}

// NewHandler creates an inventory handler backed by service.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetBalances returns balances matching optional warehouse and variant filters.
func (h *Handler) GetBalances(w http.ResponseWriter, r *http.Request) {
	warehouseID := optionalQueryID(r, "warehouse_id")
	variantID := optionalQueryID(r, "variant_id")

	balances, err := h.service.GetBalances(r.Context(), warehouseID, variantID)
	if err != nil {
		respondInventoryError(w, err)
		return
	}

	response := make([]InventoryBalanceResponse, len(balances))
	for i, balance := range balances {
		response[i] = InventoryBalanceResponse(balance)
	}

	httpx.RespondJSON(w, http.StatusOK, response)
}

// GetBalance returns one warehouse and variant balance.
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	input := inventoryBalancePathInput{
		WarehouseID: strings.TrimSpace(chi.URLParam(r, "warehouse_id")),
		VariantID:   strings.TrimSpace(chi.URLParam(r, "variant_id")),
	}
	if err := httpx.Validate(r.Context(), input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	balance, err := h.service.GetBalance(r.Context(), input.WarehouseID, input.VariantID)
	if err != nil {
		respondInventoryError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, InventoryBalanceResponse(*balance))
}

// Adjust creates a stock adjustment and its immutable movement record.
func (h *Handler) Adjust(w http.ResponseWriter, r *http.Request) {
	var input adjustmentInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	adjustment := Adjustment(input)
	balance, movement, err := h.service.Adjust(r.Context(), &adjustment)
	if err != nil {
		respondInventoryError(w, err)
		return
	}

	response := AdjustmentResponse{
		Balance:  InventoryBalanceResponse(*balance),
		Movement: StockMovementResponse(*movement),
	}
	httpx.RespondJSON(w, http.StatusCreated, response)
}

// GetMovements returns movement history matching optional warehouse and variant filters.
func (h *Handler) GetMovements(w http.ResponseWriter, r *http.Request) {
	warehouseID := optionalQueryID(r, "warehouse_id")
	variantID := optionalQueryID(r, "variant_id")

	movements, err := h.service.GetMovements(r.Context(), warehouseID, variantID)
	if err != nil {
		respondInventoryError(w, err)
		return
	}

	response := make([]StockMovementResponse, len(movements))
	for i, movement := range movements {
		response[i] = StockMovementResponse(movement)
	}

	httpx.RespondJSON(w, http.StatusOK, response)
}

// GetMovementByID returns one stock movement by ID.
func (h *Handler) GetMovementByID(w http.ResponseWriter, r *http.Request) {
	input := movementIDInput{ID: strings.TrimSpace(chi.URLParam(r, "id"))}
	if err := httpx.Validate(r.Context(), input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	movement, err := h.service.GetMovementByID(r.Context(), input.ID)
	if err != nil {
		respondInventoryError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, StockMovementResponse(*movement))
}

func optionalQueryID(r *http.Request, name string) *string {
	value := strings.TrimSpace(r.URL.Query().Get(name))
	if value == "" {
		return nil
	}

	return &value
}

func respondInventoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errInvalidAdjustmentType):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_type", "unsupported adjustment type")
	case errors.Is(err, errInvalidAdjustmentQuantity):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_quantity", "quantity must be greater than zero")
	case errors.Is(err, errInvalidWarehouseID):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_warehouse_id", "warehouse_id is required")
	case errors.Is(err, errInvalidVariantID):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_variant_id", "variant_id is required")
	case errors.Is(err, errWarehouseNotFound):
		httpx.RespondError(w, http.StatusNotFound, "warehouse_not_found", "warehouse not found")
	case errors.Is(err, errWarehouseNotActive):
		httpx.RespondError(w, http.StatusConflict, "warehouse_not_active", "warehouse is not active")
	case errors.Is(err, errVariantNotFound):
		httpx.RespondError(w, http.StatusNotFound, "variant_not_found", "product variant not found")
	case errors.Is(err, errInsufficientStock):
		httpx.RespondError(w, http.StatusConflict, "insufficient_stock", "insufficient stock")
	case errors.Is(err, sql.ErrNoRows):
		httpx.RespondError(w, http.StatusNotFound, "not_found", "inventory record not found")
	default:
		httpx.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}
