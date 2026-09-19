package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	jwtutil "bearuang-go/internal/jwt"
)

// Repository provides persistence for users.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// Password hashes passwords and verifies them against encoded hashes.
type Password interface {
	Hash(password string) (string, error)
	Verify(encodedHash, password string) bool
}

var (
	errEmailAlreadyExists  = errors.New("email already exists")
	errInvalidCredentials  = errors.New("invalid credentials")
	errInvalidRefreshToken = errors.New("invalid refresh token")
)

// service contains authentication business operations.
type service struct {
	repo      Repository
	password  Password
	jwtSecret string
}

// NewService creates an authentication service backed by repo and password.
func NewService(repo Repository, jwtSecret string, password Password) *service {
	return &service{
		repo:      repo,
		password:  password,
		jwtSecret: jwtSecret,
	}
}

// Register creates a user account without signing the user in.
func (s *service) Register(ctx context.Context, email, password string) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	passwordHash, err := s.password.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        normalizeEmail(email),
		PasswordHash: string(passwordHash),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login verifies credentials and returns an access and refresh token for the user.
func (s *service) Login(ctx context.Context, email, password string) (jwtutil.TokenPair, error) {
	user, err := s.repo.GetByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return jwtutil.TokenPair{}, errInvalidCredentials
		}

		return jwtutil.TokenPair{}, err
	}
	if user == nil || !s.password.Verify(user.PasswordHash, password) {
		return jwtutil.TokenPair{}, errInvalidCredentials
	}

	return jwtutil.GenerateTokenPair(user.ID, s.jwtSecret)
}

// Refresh validates a refresh token and returns a new access and refresh token.
func (s *service) Refresh(ctx context.Context, refreshToken string) (jwtutil.TokenPair, error) {
	if err := ctx.Err(); err != nil {
		return jwtutil.TokenPair{}, err
	}

	claims, err := jwtutil.ValidateRefreshToken(refreshToken, s.jwtSecret)
	if err != nil || claims.Subject == "" {
		return jwtutil.TokenPair{}, errInvalidRefreshToken
	}

	return jwtutil.GenerateTokenPair(claims.Subject, s.jwtSecret)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
