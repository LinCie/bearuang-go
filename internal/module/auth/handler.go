package auth

import (
	"errors"
	"net/http"

	"bearuang-go/internal/httpx"
)

// Handler handles authentication requests.
type Handler struct {
	service *Service
}

type credentialsInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// NewHandler creates an authentication handler backed by service.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register creates a user account.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input credentialsInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	user, err := h.service.Register(r.Context(), input.Email, input.Password)
	if err != nil {
		respondAuthError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusCreated, user)
}

// Login verifies a user's credentials and returns access and refresh tokens.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input credentialsInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	tokens, err := h.service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		respondAuthError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, tokens)
}

// Refresh exchanges a refresh token for a new access and refresh token.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input refreshTokenInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	tokens, err := h.service.Refresh(r.Context(), input.RefreshToken)
	if err != nil {
		respondAuthError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, tokens)
}

func respondAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errEmailAlreadyExists):
		httpx.RespondError(
			w,
			http.StatusConflict,
			"email_exists",
			"email is already registered",
		)
	case errors.Is(err, errInvalidCredentials):
		httpx.RespondError(
			w,
			http.StatusUnauthorized,
			"invalid_credentials",
			"invalid email or password",
		)
	case errors.Is(err, errInvalidRefreshToken):
		httpx.RespondError(
			w,
			http.StatusUnauthorized,
			"invalid_refresh_token",
			"invalid refresh token",
		)
	default:
		httpx.RespondError(
			w,
			http.StatusInternalServerError,
			"internal_error",
			"internal server error",
		)
	}
}
