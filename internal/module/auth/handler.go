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

type loginResponse struct {
	Token string `json:"token"`
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

// Login verifies a user's credentials and returns a JWT.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input credentialsInput
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		httpx.RespondInvalidBody(w, err)
		return
	}

	token, err := h.service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		respondAuthError(w, err)
		return
	}

	httpx.RespondJSON(w, http.StatusOK, loginResponse{Token: token})
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
	default:
		httpx.RespondError(
			w,
			http.StatusInternalServerError,
			"internal_error",
			"internal server error",
		)
	}
}
