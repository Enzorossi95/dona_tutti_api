package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents the domain entity for users
type User struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"` // Never expose password hash
	FirstName         string     `json:"first_name,omitempty"`
	LastName          string     `json:"last_name,omitempty"`
	IsActive          bool       `json:"is_active"`
	IsVerified        bool       `json:"is_verified"`
	ResetToken        *string    `json:"-"` // Never expose reset token
	ResetTokenExpires *time.Time `json:"-"`
	LastLogin         *time.Time `json:"last_login,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// AuthToken represents the authentication token response
type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // in seconds
}

// LoginCredentials represents the login request data
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PasswordReset represents the password reset request data
type PasswordReset struct {
	Email string `json:"email"`
}

// PasswordUpdate represents the password update request data
type PasswordUpdate struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ResetTokenVerification represents the reset token verification request
type ResetTokenVerification struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
