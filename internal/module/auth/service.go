package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"bearuang-go/internal/database"
	jwtutil "bearuang-go/internal/jwt"
)

var (
	errEmailAlreadyExists = errors.New("email already exists")
	errInvalidCredentials = errors.New("invalid credentials")
)

// Service contains authentication business operations.
type Service struct {
	repo      Repository
	jwtSecret string
}

// NewService creates an authentication service backed by repo.
func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a user account without signing the user in.
func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           database.GenerateULID(),
		Email:        normalizeEmail(email),
		PasswordHash: string(passwordHash),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login verifies credentials and returns a signed JWT for the user.
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errInvalidCredentials
		}

		return "", err
	}
	if user == nil || !verifyPassword(user.PasswordHash, password) {
		return "", errInvalidCredentials
	}

	return jwtutil.GenerateToken(user.ID, s.jwtSecret)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
