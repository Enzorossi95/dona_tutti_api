package user

import (
	"time"

	"github.com/google/uuid"
)

// UserModel represents the database table structure with GORM tags
type UserModel struct {
	ID                uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email             string     `gorm:"uniqueIndex;not null"`
	PasswordHash      string     `gorm:"column:password_hash;not null"`
	RoleID            uuid.UUID  `gorm:"column:role_id;not null"`
	FirstName         string     `gorm:"column:first_name"`
	LastName          string     `gorm:"column:last_name"`
	IsActive          bool       `gorm:"column:is_active;default:true"`
	IsVerified        bool       `gorm:"column:is_verified;default:false"`
	ResetToken        *string    `gorm:"column:reset_token"`
	ResetTokenExpires *time.Time `gorm:"column:reset_token_expires_at"`
	LastLogin         *time.Time `gorm:"column:last_login"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}

// ToEntity converts a database model to a domain entity
func (m UserModel) ToEntity() User {
	return User{
		ID:                m.ID,
		Email:             m.Email,
		PasswordHash:      m.PasswordHash,
		RoleID:            m.RoleID,
		FirstName:         m.FirstName,
		LastName:          m.LastName,
		IsActive:          m.IsActive,
		IsVerified:        m.IsVerified,
		ResetToken:        m.ResetToken,
		ResetTokenExpires: m.ResetTokenExpires,
		LastLogin:         m.LastLogin,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *UserModel) FromEntity(entity User) {
	m.ID = entity.ID
	m.Email = entity.Email
	m.PasswordHash = entity.PasswordHash
	m.RoleID = entity.RoleID
	m.FirstName = entity.FirstName
	m.LastName = entity.LastName
	m.IsActive = entity.IsActive
	m.IsVerified = entity.IsVerified
	m.ResetToken = entity.ResetToken
	m.ResetTokenExpires = entity.ResetTokenExpires
	m.LastLogin = entity.LastLogin
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}
