package productcategories

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"bearuang-go/internal/httpx"
)

/*
======================================
Product Categories
======================================
*/

// Service provides product category business operations to the handler.
type Service interface {
	Create(ctx context.Context, category *ProductCategory) error
	GetByID(ctx context.Context, id string) (*ProductCategory, error)
	GetMany(ctx context.Context) ([]ProductCategory, error)
	Update(ctx context.Context, category *ProductCategory) error
	Delete(ctx context.Context, id string) error
}

// Handler handles HTTP requests for product categories.
type Handler struct {
	service Service
}

/*
======================================
Product Category Types
======================================
*/

type productCategoryCreateInput struct {
	ParentID    *string `json:"parent_id"`
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status" validate:"required,oneof=draft active inactive archived"`
}

type productCategoryIDInput struct {
	ID string `validate:"required"`
}

/*
======================================
Product Category Handlers
======================================
*/

// NewHandler creates a handler for product categories backed by service.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetMany returns all non-deleted product categories.
func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetMany(r.Context())
	if err != nil {
		respondProductCategoryError(w, err)
		return
	}

	response := make([]ProductCategoryResponse, len(categories))
	for i, c := range categories {
		response[i] = ProductCategoryResponse(c)
	}

	httpx.RespondJSON(w, http.StatusOK, response)
}

// GetByID returns a non-deleted product category by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := validateProductCategoryID(w, r)
	if !ok {
		return
	}

	category, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondProductCategoryError(w, err)
		return
	}

	response := ProductCategoryResponse(*category)

	httpx.RespondJSON(w, http.StatusOK, response)
}

// Create creates a product category.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input productCategoryCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	category := ProductCategory{
		ParentID:    input.ParentID,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Status:      input.Status,
	}
	if err := h.service.Create(r.Context(), &category); err != nil {
		respondProductCategoryError(w, err)
		return
	}

	response := ProductCategoryResponse(category)

	httpx.RespondJSON(w, http.StatusCreated, response)
}

// Update updates a non-deleted product category by ID.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := validateProductCategoryID(w, r)
	if !ok {
		return
	}

	var input productCategoryCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	category := ProductCategory{
		ID:          id,
		ParentID:    input.ParentID,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Status:      input.Status,
	}
	if err := h.service.Update(r.Context(), &category); err != nil {
		respondProductCategoryError(w, err)
		return
	}

	response := ProductCategoryResponse(category)

	httpx.RespondJSON(w, http.StatusOK, response)
}

// Delete soft-deletes a non-deleted product category by ID.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := validateProductCategoryID(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondProductCategoryError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, nil)
}

/*
======================================
Product Category Helpers
======================================
*/

func validateProductCategoryID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if err := httpx.Validate(r.Context(), productCategoryIDInput{ID: id}); err != nil {
		httpx.RespondInvalidBody(w, err)
		return "", false
	}

	return id, true
}

func respondProductCategoryError(w http.ResponseWriter, err error) {
	if errors.Is(err, errInvalidProductCategoryStatus) {
		httpx.RespondError(w, http.StatusBadRequest, "invalid_status", err.Error())
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		httpx.RespondError(w, http.StatusNotFound, "not_found", "product category not found")
		return
	}

	httpx.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
}
