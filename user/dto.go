package user

import "github.com/google/uuid"

// RegisterDTO represents the request body for user registration
type RegisterDTO struct {
	Email     string `json:"email" example:"user@example.com"`
	Password  string `json:"password" example:"strongpassword123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// RegisterResponseDTO represents the response for user registration
type RegisterResponseDTO struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// LoginDTO represents the request body for user login
type LoginDTO struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongpassword123"`
}

// LoginResponseDTO represents the response for user login
type LoginResponseDTO struct {
	Token AuthToken `json:"token"`
}

// GetUserDTO represents the request parameters for getting a user
type GetUserDTO struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// GetUserResponseDTO represents the response for getting a user
type GetUserResponseDTO struct {
	User User `json:"user"`
}

// UpdateUserDTO represents the request body for updating a user
type UpdateUserDTO struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// UpdatePasswordDTO represents the request body for updating a user's password
type UpdatePasswordDTO struct {
	CurrentPassword string `json:"current_password" example:"oldpassword123"`
	NewPassword     string `json:"new_password" example:"newpassword123"`
}

// RequestPasswordResetDTO represents the request body for requesting a password reset
type RequestPasswordResetDTO struct {
	Email string `json:"email" example:"user@example.com"`
}

// ResetPasswordDTO represents the request body for resetting a password
type ResetPasswordDTO struct {
	Token       string `json:"token" example:"reset-token-123"`
	NewPassword string `json:"new_password" example:"newpassword123"`
}
