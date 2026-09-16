package products

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"bearuang-go/internal/httpx"
)

/*
======================================
Products
======================================
*/

// Handler handles HTTP requests for products and their variants.
type Handler struct {
	service *Service
}

type productWriteInput struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"required,oneof=draft active inactive archived"`
}

type productCreateInput struct {
	productWriteInput
	ID string `json:"id" validate:"required"`
}

type productIDInput struct {
	ID string `validate:"required"`
}

// NewHandler creates a handler for products and their variants backed by service.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetMany returns all non-deleted products.
func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetMany(r.Context())
	if err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productDTOs(products))
}

// GetByID returns a non-deleted product by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := validateProductID(w, r)
	if !ok {
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
	var input productCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	product, err := h.service.Create(r.Context(), &Product{
		ID:          input.ID,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Status:      input.Status,
	})
	if err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusCreated, productDTO(product))
}

// Update updates a non-deleted product by ID.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := validateProductID(w, r)
	if !ok {
		return
	}

	var input productWriteInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	product, err := h.service.Update(r.Context(), &Product{
		ID:          id,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Status:      input.Status,
	})
	if err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productDTO(product))
}

// Delete soft-deletes a non-deleted product by ID.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := validateProductID(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondProductError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, nil)
}

func validateProductID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := strings.TrimSpace(r.PathValue("id"))
	if err := httpx.Validate(r.Context(), productIDInput{ID: id}); err != nil {
		httpx.RespondInvalidBody(w, err)
		return "", false
	}

	return id, true
}

func respondProductError(w http.ResponseWriter, err error) {
	if errors.Is(err, errInvalidProductStatus) {
		httpx.RespondError(w, http.StatusBadRequest, "invalid_status", err.Error())
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		httpx.RespondError(w, http.StatusNotFound, "not_found", "product not found")
		return
	}

	httpx.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
}

/*
======================================
Variants
======================================
*/

type variantWriteInput struct {
	SKU    string  `json:"sku" validate:"required"`
	Name   string  `json:"name" validate:"required"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
	Unit   string  `json:"unit" validate:"required"`
	Status string  `json:"status" validate:"required,oneof=draft active inactive archived"`
}

type variantCreateInput struct {
	variantWriteInput
	ID string `json:"id" validate:"required"`
}

type variantProductIDInput struct {
	ProductID string `json:"product_id" validate:"required"`
}

type variantPathInput struct {
	ProductID string `json:"product_id" validate:"required"`
	VariantID string `json:"variant_id" validate:"required"`
}

// GetManyVariants returns all non-deleted variants for a product.
func (h *Handler) GetManyVariants(w http.ResponseWriter, r *http.Request) {
	productID, ok := validateVariantProductID(w, r)
	if !ok {
		return
	}

	variants, err := h.service.GetManyVariantsByProduct(r.Context(), productID)
	if err != nil {
		respondVariantError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productVariantDTOs(variants))
}

// GetVariantByID returns a non-deleted product variant by ID.
func (h *Handler) GetVariantByID(w http.ResponseWriter, r *http.Request) {
	productID, variantID, ok := validateVariantPath(w, r)
	if !ok {
		return
	}

	variant, err := h.service.GetVariantByID(r.Context(), productID, variantID)
	if err != nil {
		respondVariantError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productVariantDTO(variant))
}

// CreateVariant creates a product variant.
func (h *Handler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	productID, ok := validateVariantProductID(w, r)
	if !ok {
		return
	}

	var input variantCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	variant, err := h.service.CreateVariant(r.Context(), &ProductVariant{
		ID:        input.ID,
		ProductID: productID,
		SKU:       input.SKU,
		Name:      input.Name,
		Price:     input.Price,
		Stock:     input.Stock,
		Unit:      input.Unit,
		Status:    input.Status,
	})
	if err != nil {
		respondVariantError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusCreated, productVariantDTO(variant))
}

// UpdateVariant updates a non-deleted product variant by ID.
func (h *Handler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	productID, variantID, ok := validateVariantPath(w, r)
	if !ok {
		return
	}

	var input variantWriteInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	variant, err := h.service.UpdateVariant(r.Context(), &ProductVariant{
		ID:        variantID,
		ProductID: productID,
		SKU:       input.SKU,
		Name:      input.Name,
		Price:     input.Price,
		Stock:     input.Stock,
		Unit:      input.Unit,
		Status:    input.Status,
	})
	if err != nil {
		respondVariantError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, productVariantDTO(variant))
}

// DeleteVariant soft-deletes a non-deleted product variant by ID.
func (h *Handler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	productID, variantID, ok := validateVariantPath(w, r)
	if !ok {
		return
	}

	if err := h.service.DeleteVariant(r.Context(), productID, variantID); err != nil {
		respondVariantError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, nil)
}

func validateVariantProductID(w http.ResponseWriter, r *http.Request) (string, bool) {
	productID := strings.TrimSpace(r.PathValue("product_id"))
	if err := httpx.Validate(r.Context(), variantProductIDInput{ProductID: productID}); err != nil {
		httpx.RespondInvalidBody(w, err)
		return "", false
	}

	return productID, true
}

func validateVariantPath(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	input := variantPathInput{
		ProductID: strings.TrimSpace(r.PathValue("product_id")),
		VariantID: strings.TrimSpace(r.PathValue("variant_id")),
	}
	if err := httpx.Validate(r.Context(), input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return "", "", false
	}

	return input.ProductID, input.VariantID, true
}

func respondVariantError(w http.ResponseWriter, err error) {
	if errors.Is(err, errInvalidVariantStatus) {
		httpx.RespondError(w, http.StatusBadRequest, "invalid_status", err.Error())
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		httpx.RespondError(w, http.StatusNotFound, "not_found", "product variant not found")
		return
	}

	httpx.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
}
