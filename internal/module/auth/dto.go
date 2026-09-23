package auth

import "time"

/*
======================================
Users
======================================
*/

// UserResponse represents a user response.
type UserResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
