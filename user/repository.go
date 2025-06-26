package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByEmailWithRole(ctx context.Context, email string) (User, string, error)
	GetUserByResetToken(ctx context.Context, token string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	SetResetToken(ctx context.Context, id uuid.UUID, token string, expires time.Time) error
	ClearResetToken(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return model.ToEntity(), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return model.ToEntity(), nil
}

func (r *userRepository) GetUserByResetToken(ctx context.Context, token string) (User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).
		Where("reset_token = ? AND reset_token_expires_at > ?", token, time.Now()).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return User{}, fmt.Errorf("invalid or expired reset token")
		}
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return model.ToEntity(), nil
}

func (r *userRepository) ListUsers(ctx context.Context) ([]User, error) {
	var models []UserModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}
	return users, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user User) error {
	model := UserModel{}
	model.FromEntity(user)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user User) error {
	model := UserModel{}
	model.FromEntity(user)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).
		Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *userRepository) SetResetToken(ctx context.Context, id uuid.UUID, token string, expires time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"reset_token":            token,
			"reset_token_expires_at": expires,
		}).Error; err != nil {
		return fmt.Errorf("failed to set reset token: %w", err)
	}
	return nil
}

func (r *userRepository) ClearResetToken(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"reset_token":            nil,
			"reset_token_expires_at": nil,
		}).Error; err != nil {
		return fmt.Errorf("failed to clear reset token: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", id).
		Update("last_login", time.Now()).
		Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

func (r *userRepository) GetUserByEmailWithRole(ctx context.Context, email string) (User, string, error) {
	var result struct {
		UserModel
		RoleName string `gorm:"column:role_name"`
	}

	query := `
		SELECT u.*, r.name as role_name
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.email = ? AND r.is_active = true
	`

	if err := r.db.WithContext(ctx).Raw(query, email).Scan(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return User{}, "", fmt.Errorf("user not found")
		}
		return User{}, "", fmt.Errorf("failed to get user with role: %w", err)
	}

	return result.UserModel.ToEntity(), result.RoleName, nil
}
