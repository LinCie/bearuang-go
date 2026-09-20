package middleware

import (
	"net/http"
	"strings"

	"bearuang-go/internal/config"
	"bearuang-go/internal/httpx"
	jwtutil "bearuang-go/internal/jwt"
)

// Auth requires a valid access token cookie or Bearer token before serving the next handler.
func Auth(cfg config.Config) httpx.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := accessToken(r)
			if token == "" {
				respondUnauthorized(w)
				return
			}

			if _, err := jwtutil.ValidateAccessToken(token, cfg.JWTSecret); err != nil {
				respondUnauthorized(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func accessToken(r *http.Request) string {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}

	cookie, err := r.Cookie(config.AccessTokenCookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func respondUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	httpx.RespondError(
		w,
		http.StatusUnauthorized,
		"unauthorized",
		"authentication required",
	)
}
