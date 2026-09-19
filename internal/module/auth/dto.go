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

// toUserResponse converts a user to its response DTO.
func toUserResponse(user User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
