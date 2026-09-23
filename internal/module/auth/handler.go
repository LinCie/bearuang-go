package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"bearuang-go/internal/config"
	"bearuang-go/internal/httpx"
	jwtutil "bearuang-go/internal/jwt"
)

/*
======================================
Authentication
======================================
*/

// Service provides authentication business operations to the handler.
type Service interface {
	Register(ctx context.Context, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (jwtutil.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (jwtutil.TokenPair, error)
}

// Handler handles authentication requests.
type Handler struct {
	service      Service
	cookieSecure bool
}

/*
======================================
Authentication Types
======================================
*/

type credentialsInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type authSessionResponse struct {
	Authenticated bool `json:"authenticated"`
}

/*
======================================
Authentication Handlers
======================================
*/

// NewHandler creates an authentication handler backed by service.
func NewHandler(service Service, cookieSecure bool) *Handler {
	return &Handler{
		service:      service,
		cookieSecure: cookieSecure,
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

	httpx.RespondJSON(w, http.StatusCreated, UserResponse(*user))
}

// Login verifies credentials and sets HTTP-only access and refresh cookies.
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

	h.setTokenCookies(w, tokens)
	respondAuthenticated(w)
}

// Refresh exchanges the refresh cookie for new access and refresh cookies.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie(config.RefreshTokenCookieName)
	if err != nil {
		respondAuthError(w, errInvalidRefreshToken)
		return
	}

	tokens, err := h.service.Refresh(r.Context(), refreshCookie.Value)
	if err != nil {
		respondAuthError(w, err)
		return
	}

	h.setTokenCookies(w, tokens)
	respondAuthenticated(w)
}

// Logout clears the authentication cookies.
func (h *Handler) Logout(w http.ResponseWriter, _ *http.Request) {
	h.clearTokenCookies(w)
	respondSession(w, false)
}

/*
======================================
Authentication Helpers
======================================
*/

func (h *Handler) setTokenCookies(w http.ResponseWriter, tokens jwtutil.TokenPair) {
	http.SetCookie(w, newTokenCookie(
		config.AccessTokenCookieName,
		tokens.AccessToken,
		jwtutil.AccessTokenLifetime,
		h.cookieSecure,
	))
	http.SetCookie(w, newTokenCookie(
		config.RefreshTokenCookieName,
		tokens.RefreshToken,
		jwtutil.RefreshTokenLifetime,
		h.cookieSecure,
	))
}

func (h *Handler) clearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, expiredTokenCookie(config.AccessTokenCookieName, h.cookieSecure))
	http.SetCookie(w, expiredTokenCookie(config.RefreshTokenCookieName, h.cookieSecure))
}

func newTokenCookie(name, value string, lifetime time.Duration, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(lifetime / time.Second),
		Expires:  time.Now().Add(lifetime),
	}
}

func expiredTokenCookie(name string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
	}
}

func respondAuthenticated(w http.ResponseWriter) {
	respondSession(w, true)
}

func respondSession(w http.ResponseWriter, authenticated bool) {
	w.Header().Set("Cache-Control", "no-store")
	httpx.RespondJSON(w, http.StatusOK, authSessionResponse{Authenticated: authenticated})
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
