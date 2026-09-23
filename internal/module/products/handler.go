package products

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
Products
======================================
*/

// Service provides product and product variant business operations to the handler.
type Service interface {
	/*
		======================================
		Products
		======================================
	*/

	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	GetMany(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error

	/*
		======================================
		Variants
		======================================
	*/

	CreateVariant(ctx context.Context, variant *ProductVariant) error
	GetVariantByID(ctx context.Context, productID, id string) (*ProductVariant, error)
	GetManyVariantsByProduct(ctx context.Context, productID string) ([]ProductVariant, error)
	UpdateVariant(ctx context.Context, variant *ProductVariant) error
	DeleteVariant(ctx context.Context, productID, id string) error
}

// Handler handles HTTP requests for products and their variants.
type Handler struct {
	service Service
}

/*
======================================
Product Types
======================================
*/

type productCreateInput struct {
	CategoryID  *string `json:"category_id"`
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status" validate:"required,oneof=draft active inactive archived"`
}

type productIDInput struct {
	ID string `validate:"required"`
}

/*
======================================
Product Handlers
======================================
*/

// NewHandler creates a handler for products and their variants backed by service.
func NewHandler(service Service) *Handler {
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

	response := make([]ProductResponse, len(products))
	for i, product := range products {
		response[i] = ProductResponse(product)
	}

	httpx.RespondJSON(w, http.StatusOK, response)
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

	response := ProductResponse(*product)

	httpx.RespondJSON(w, http.StatusOK, response)
}

// Create creates a product.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input productCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	product := Product{
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Status:      input.Status,
	}
	err := h.service.Create(r.Context(), &product)
	if err != nil {
		respondProductError(w, err)
		return
	}

	response := ProductResponse(product)

	httpx.RespondJSON(w, http.StatusCreated, response)
}

// Update updates a non-deleted product by ID.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := validateProductID(w, r)
	if !ok {
		return
	}

	var input productCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	product := Product{
		ID:          id,
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Status:      input.Status,
	}
	if err := h.service.Update(r.Context(), &product); err != nil {
		respondProductError(w, err)
		return
	}

	response := ProductResponse(product)

	httpx.RespondJSON(w, http.StatusOK, response)
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

/*
======================================
Product Helpers
======================================
*/

func validateProductID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
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
Variant Types
======================================
*/

type variantCreateInput struct {
	SKU    string  `json:"sku" validate:"required"`
	Name   string  `json:"name" validate:"required"`
	Price  float64 `json:"price"`
	Unit   string  `json:"unit" validate:"required"`
	Status string  `json:"status" validate:"required,oneof=draft active inactive archived"`
}

type variantProductIDInput struct {
	ProductID string `json:"product_id" validate:"required"`
}

type variantPathInput struct {
	ProductID string `json:"product_id" validate:"required"`
	VariantID string `json:"variant_id" validate:"required"`
}

/*
======================================
Variant Handlers
======================================
*/

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

	response := make([]ProductVariantResponse, len(variants))
	for i, variant := range variants {
		response[i] = ProductVariantResponse(variant)
	}

	httpx.RespondJSON(w, http.StatusOK, response)
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

	response := ProductVariantResponse(*variant)

	httpx.RespondJSON(w, http.StatusOK, response)
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

	variant := ProductVariant{
		ProductID: productID,
		SKU:       input.SKU,
		Name:      input.Name,
		Price:     input.Price,
		Unit:      input.Unit,
		Status:    input.Status,
	}
	if err := h.service.CreateVariant(r.Context(), &variant); err != nil {
		respondVariantError(w, err)
		return
	}

	response := ProductVariantResponse(variant)

	httpx.RespondJSON(w, http.StatusCreated, response)
}

// UpdateVariant updates a non-deleted product variant by ID.
func (h *Handler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	productID, variantID, ok := validateVariantPath(w, r)
	if !ok {
		return
	}

	var input variantCreateInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	variant := ProductVariant{
		ID:        variantID,
		ProductID: productID,
		SKU:       input.SKU,
		Name:      input.Name,
		Price:     input.Price,
		Unit:      input.Unit,
		Status:    input.Status,
	}
	if err := h.service.UpdateVariant(r.Context(), &variant); err != nil {
		respondVariantError(w, err)
		return
	}

	response := ProductVariantResponse(variant)

	httpx.RespondJSON(w, http.StatusOK, response)
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

/*
======================================
Variant Helpers
======================================
*/

func validateVariantProductID(w http.ResponseWriter, r *http.Request) (string, bool) {
	productID := strings.TrimSpace(chi.URLParam(r, "product_id"))
	if err := httpx.Validate(r.Context(), variantProductIDInput{ProductID: productID}); err != nil {
		httpx.RespondInvalidBody(w, err)
		return "", false
	}

	return productID, true
}

func validateVariantPath(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	input := variantPathInput{
		ProductID: strings.TrimSpace(chi.URLParam(r, "product_id")),
		VariantID: strings.TrimSpace(chi.URLParam(r, "variant_id")),
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
