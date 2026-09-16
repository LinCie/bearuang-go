package middleware

import (
	"net/http"
	"strings"

	"bearuang-go/internal/httpx"
	jwtutil "bearuang-go/internal/jwt"
)

// Auth requires a valid Bearer access token before serving the next handler.
func Auth(secret string) httpx.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Fields(r.Header.Get("Authorization"))
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				respondUnauthorized(w)
				return
			}

			if _, err := jwtutil.ValidateAccessToken(parts[1], secret); err != nil {
				respondUnauthorized(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
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
