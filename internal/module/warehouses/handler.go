package warehouses

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"bearuang-go/internal/httpx"
)

// Service provides warehouse business operations to the handler.
type Service interface {
	Create(ctx context.Context, warehouse *Warehouse) error
	GetByID(ctx context.Context, id string) (*Warehouse, error)
	GetMany(ctx context.Context) ([]Warehouse, error)
	Update(ctx context.Context, warehouse *Warehouse) error
	Delete(ctx context.Context, id string) error
}

// Handler handles HTTP requests for warehouses.
type Handler struct {
	service Service
}

type warehouseInput struct {
	Code        string  `json:"code" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Address     *string `json:"address"`
	Status      string  `json:"status" validate:"required"`
}

type warehouseIDInput struct {
	ID string `validate:"required"`
}

// NewHandler creates a handler for warehouses backed by service.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetMany returns all non-deleted warehouses.
func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.service.GetMany(r.Context())
	if err != nil {
		respondWarehouseError(w, err)
		return
	}

	response := make([]WarehouseResponse, len(warehouses))
	for i, warehouse := range warehouses {
		response[i] = WarehouseResponse(warehouse)
	}

	httpx.RespondJSON(w, http.StatusOK, response)
}

// GetByID returns a non-deleted warehouse by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	warehouse, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondWarehouseError(w, err)
		return
	}

	response := WarehouseResponse(*warehouse)
	httpx.RespondJSON(w, http.StatusOK, response)
}

// Create creates a warehouse.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input warehouseInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	warehouse := Warehouse{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		Address:     input.Address,
		Status:      input.Status,
	}
	if err := h.service.Create(r.Context(), &warehouse); err != nil {
		respondWarehouseError(w, err)
		return
	}

	response := WarehouseResponse(warehouse)
	httpx.RespondJSON(w, http.StatusCreated, response)
}

// Update updates a non-deleted warehouse by ID.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	var input warehouseInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	warehouse := Warehouse{
		ID:          id,
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		Address:     input.Address,
		Status:      input.Status,
	}
	if err := h.service.Update(r.Context(), &warehouse); err != nil {
		respondWarehouseError(w, err)
		return
	}

	response := WarehouseResponse(warehouse)
	httpx.RespondJSON(w, http.StatusOK, response)
}

// Delete soft-deletes a non-deleted warehouse by ID.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := validateWarehouseID(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondWarehouseError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, nil)
}

func validateWarehouseID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if err := httpx.Validate(r.Context(), warehouseIDInput{ID: id}); err != nil {
		httpx.RespondInvalidBody(w, err)
		return "", false
	}

	return id, true
}

func respondWarehouseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errInvalidWarehouseCode):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_code", err.Error())
	case errors.Is(err, errInvalidWarehouseName):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_name", err.Error())
	case errors.Is(err, errInvalidWarehouseStatus):
		httpx.RespondError(w, http.StatusBadRequest, "invalid_status", err.Error())
	case errors.Is(err, errDuplicateWarehouseCode):
		httpx.RespondError(w, http.StatusConflict, "duplicate_code", err.Error())
	case errors.Is(err, errWarehouseInUse):
		httpx.RespondError(w, http.StatusConflict, "warehouse_in_use", err.Error())
	case errors.Is(err, errWarehouseHasStock):
		httpx.RespondError(w, http.StatusConflict, "warehouse_has_stock", err.Error())
	case errors.Is(err, sql.ErrNoRows):
		httpx.RespondError(w, http.StatusNotFound, "not_found", "warehouse not found")
	default:
		httpx.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}
