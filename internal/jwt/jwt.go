package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const tokenLifetime = 24 * time.Hour

var (
	// ErrSecretRequired is returned when no JWT signing secret is provided.
	ErrSecretRequired = errors.New("jwt secret is required")
	// ErrInvalidToken is returned when a token cannot be validated.
	ErrInvalidToken = errors.New("invalid jwt token")
)

// Claims contains the claims issued by this package.
type Claims struct {
	jwtv5.RegisteredClaims
}

// GenerateToken creates an HS256 JWT for subject that expires after 24 hours.
func GenerateToken(subject, secret string) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrSecretRequired
	}

	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(tokenLifetime)),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signedToken, nil
}

// ValidateToken verifies the signature and claims of a JWT and returns its claims.
func ValidateToken(tokenString, secret string) (*Claims, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, ErrSecretRequired
	}
	if strings.TrimSpace(tokenString) == "" {
		return nil, ErrInvalidToken
	}

	claims := &Claims{}
	token, err := jwtv5.NewParser(
		jwtv5.WithValidMethods([]string{jwtv5.SigningMethodHS256.Alg()}),
		jwtv5.WithExpirationRequired(),
	).ParseWithClaims(tokenString, claims, func(token *jwtv5.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
