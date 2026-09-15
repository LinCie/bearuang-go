package products

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"bearuang-go/internal/httpx"
)

// Handler handles HTTP requests for products.
type Handler struct {
	service *Service
}

type productInput struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"required"`
	ID          string `json:"id" validate:"required"`
}

type productIDInput struct {
	ID string `validate:"required"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

// NewHandler creates a product handler backed by service.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetMany returns all active products.
func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetMany(r.Context())
	if err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productDTOs(products))
}

// GetByID returns an active product by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if err := validate.Struct(productIDInput{ID: id}); err != nil {
		httpx.RespondValidationError(w, err)
		return
	}

	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productDTO(product))
}

// Create creates a product.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input productInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	product, err := h.service.Create(r.Context(), input.product(input.ID))
	if err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusCreated, productDTO(product))
}

// Update updates an active product by ID.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if err := validate.Struct(productIDInput{ID: id}); err != nil {
		httpx.RespondValidationError(w, err)
		return
	}

	input := productInput{ID: id}
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	product, err := h.service.Update(r.Context(), input.product(id))
	if err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productDTO(product))
}

// Delete soft-deletes an active product by ID.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if err := validate.Struct(productIDInput{ID: id}); err != nil {
		httpx.RespondValidationError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, nil)
}

func (p productInput) product(id string) *Product {
	return &Product{
		ID:          id,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Status:      p.Status,
	}
}

func respondProductError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		httpx.RespondError(w, http.StatusNotFound, "not_found", "product not found")
		return
	}

	httpx.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
}
