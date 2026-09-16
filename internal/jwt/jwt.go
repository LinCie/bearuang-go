package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const (
	// AccessTokenLifetime is the lifetime of an access token.
	AccessTokenLifetime = 15 * time.Minute
	// RefreshTokenLifetime is the lifetime of a refresh token.
	RefreshTokenLifetime = 7 * 24 * time.Hour

	// AccessTokenType identifies access tokens in their claims.
	AccessTokenType = "access"
	// RefreshTokenType identifies refresh tokens in their claims.
	RefreshTokenType = "refresh"
)

var (
	// ErrSecretRequired is returned when no JWT signing secret is provided.
	ErrSecretRequired = errors.New("jwt secret is required")
	// ErrInvalidToken is returned when a token cannot be validated.
	ErrInvalidToken = errors.New("invalid jwt token")
)

// Claims contains the claims issued by this package.
type Claims struct {
	TokenType string `json:"token_type"`
	jwtv5.RegisteredClaims
}

// TokenPair contains an access token and a refresh token.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateAccessToken creates an HS256 access token for subject.
func GenerateAccessToken(subject, secret string) (string, error) {
	return generateToken(subject, secret, AccessTokenType, AccessTokenLifetime)
}

// GenerateRefreshToken creates an HS256 refresh token for subject.
func GenerateRefreshToken(subject, secret string) (string, error) {
	return generateToken(subject, secret, RefreshTokenType, RefreshTokenLifetime)
}

// GenerateTokenPair creates an access token and a refresh token for subject.
func GenerateTokenPair(subject, secret string) (TokenPair, error) {
	accessToken, err := GenerateAccessToken(subject, secret)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := GenerateRefreshToken(subject, secret)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateToken(subject, secret, tokenType string, lifetime time.Duration) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrSecretRequired
	}

	now := time.Now()
	claims := Claims{
		TokenType: tokenType,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(lifetime)),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signedToken, nil
}

// ValidateAccessToken verifies an access token and returns its claims.
func ValidateAccessToken(tokenString, secret string) (*Claims, error) {
	return validateToken(tokenString, secret, AccessTokenType)
}

// ValidateRefreshToken verifies a refresh token and returns its claims.
func ValidateRefreshToken(tokenString, secret string) (*Claims, error) {
	return validateToken(tokenString, secret, RefreshTokenType)
}

func validateToken(tokenString, secret, expectedType string) (*Claims, error) {
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
	if token == nil || !token.Valid || claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateToken creates an access token for subject.
//
// Deprecated: use GenerateAccessToken instead.
func GenerateToken(subject, secret string) (string, error) {
	return GenerateAccessToken(subject, secret)
}

// ValidateToken verifies an access token and returns its claims.
//
// Deprecated: use ValidateAccessToken instead.
func ValidateToken(tokenString, secret string) (*Claims, error) {
	return ValidateAccessToken(tokenString, secret)
}
