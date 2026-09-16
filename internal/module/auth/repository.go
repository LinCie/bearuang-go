package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

// Repository provides persistence for users.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// PostgresRepository stores users in PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a repository for users backed by PostgreSQL.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

var _ Repository = (*PostgresRepository)(nil)

// Create inserts a user and populates its stored fields.
func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `
		INSERT INTO USERS (
			id,
			email,
			password_hash
		) VALUES (:id, :email, :password_hash)
		RETURNING
			id,
			email,
			password_hash,
			created_at,
			updated_at`
	query, args, err := r.db.BindNamed(query, user)
	if err != nil {
		return err
	}

	if err := r.db.GetContext(ctx, user, query, args...); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return errEmailAlreadyExists
		}

		return err
	}

	return nil
}

// GetByEmail returns a user by email.
func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id,
			email,
			password_hash,
			created_at,
			updated_at
		FROM USERS
		WHERE email = $1`

	user := new(User)
	if err := r.db.GetContext(ctx, user, query, email); err != nil {
		return nil, err
	}

	return user, nil
}
