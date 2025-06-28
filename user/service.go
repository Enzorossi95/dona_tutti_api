package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	apierrors "dona_tutti_api/errors"
	"dona_tutti_api/organizer"
	"dona_tutti_api/rbac"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenExpiration      = 24 * time.Hour
	resetTokenExpiration = 1 * time.Hour
	jwtSecret            = "your-secret-key"
)

type Service interface {
	Register(ctx context.Context, email, password, firstName, lastName string) (uuid.UUID, error)
	Login(ctx context.Context, email, password string) (*AuthToken, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetMe(ctx context.Context, id uuid.UUID) (MeResponseDTO, error)
	ListUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, dto RegisterDTO) (uuid.UUID, error)
	UpdateUser(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) error
	UpdatePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type service struct {
	repo             UserRepository
	organizerService organizer.Service
}

func NewService(repo UserRepository, organizerService organizer.Service) Service {
	return &service{
		repo:             repo,
		organizerService: organizerService,
	}
}

func (s *service) Register(ctx context.Context, email, password, firstName, lastName string) (uuid.UUID, error) {
	// Validate email format
	if !isValidEmail(email) {
		return uuid.Nil, apierrors.NewFieldValidationError("email", "invalid email format")
	}

	// Check if user already exists
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return uuid.Nil, apierrors.NewFieldValidationError("email", "email already registered")
	}

	// Validate password strength
	if !isStrongPassword(password) {
		return uuid.Nil, apierrors.NewFieldValidationError("password", "password must be at least 8 characters long and contain letters and numbers")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       rbac.GuestRoleID, // Default role for new users
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create associated organizer using injected service
	organizerEntity := organizer.Organizer{
		ID:       uuid.New(),
		UserID:   user.ID,
		Name:     user.FirstName + " " + user.LastName,
		Avatar:   "",
		Verified: false,
	}
	fmt.Println("organizerEntity", organizerEntity)

	_, err = s.organizerService.CreateOrganizer(ctx, organizerEntity)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create organizer: %w", err)
	}

	return user.ID, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*AuthToken, error) {
	user, roleName, err := s.repo.GetUserByEmailWithRole(ctx, email)
	if err != nil {
		return nil, apierrors.NewValidationError("invalid email or password")
	}

	if !user.IsActive {
		return nil, apierrors.NewValidationError("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apierrors.NewValidationError("invalid email or password")
	}

	// Generate JWT token with role information
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.ID.String(),
		"role_id": user.RoleID.String(),
		"role":    roleName,
		"exp":     time.Now().Add(tokenExpiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("failed to update last login: %v\n", err)
	}

	return &AuthToken{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int(tokenExpiration.Seconds()),
	}, nil
}

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) GetMe(ctx context.Context, id uuid.UUID) (MeResponseDTO, error) {
	user, roleName, roleID, err := s.repo.GetUserByIDWithRole(ctx, id)
	if err != nil {
		return MeResponseDTO{}, err
	}

	response := MeResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role: RoleInfo{
			ID:   roleID,
			Name: roleName,
		},
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

func (s *service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *service) CreateUser(ctx context.Context, dto RegisterDTO) (uuid.UUID, error) {
	return s.Register(ctx, dto.Email, dto.Password, dto.FirstName, dto.LastName)
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.UpdatedAt = time.Now()

	return s.repo.UpdateUser(ctx, user)
}

func (s *service) UpdatePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return apierrors.NewValidationError("current password is incorrect")
	}

	if !isStrongPassword(newPassword) {
		return apierrors.NewFieldValidationError("new_password", "password must be at least 8 characters long and contain letters and numbers")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.UpdatePassword(ctx, id, string(hashedPassword))
}

func (s *service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists
		return nil
	}

	// Generate reset token
	token, err := generateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expires := time.Now().Add(resetTokenExpiration)
	if err := s.repo.SetResetToken(ctx, user.ID, token, expires); err != nil {
		return fmt.Errorf("failed to set reset token: %w", err)
	}

	// TODO: Send reset email with token
	fmt.Printf("Reset token for %s: %s\n", email, token)

	return nil
}

func (s *service) ResetPassword(ctx context.Context, token, newPassword string) error {
	if !isStrongPassword(newPassword) {
		return apierrors.NewFieldValidationError("password", "password must be at least 8 characters long and contain letters and numbers")
	}

	user, err := s.repo.GetUserByResetToken(ctx, token)
	if err != nil {
		return apierrors.NewValidationError("invalid or expired reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.UpdatePassword(ctx, user.ID, string(hashedPassword)); err != nil {
		return err
	}

	return s.repo.ClearResetToken(ctx, user.ID)
}

// Helper functions

func generateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func isValidEmail(email string) bool {
	// TODO: Implement proper email validation
	return len(email) > 0 && len(email) <= 255
}

func isStrongPassword(password string) bool {
	// TODO: Implement proper password strength validation
	return len(password) >= 8
}
