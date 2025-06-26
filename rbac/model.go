package rbac

import (
	"time"

	"github.com/google/uuid"
)

// RoleModel represents the database table structure for roles
type RoleModel struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string    `gorm:"column:name;uniqueIndex;not null"`
	Description string    `gorm:"column:description"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (RoleModel) TableName() string {
	return "roles"
}

// ToEntity converts a database model to a domain entity
func (m RoleModel) ToEntity() Role {
	return Role{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *RoleModel) FromEntity(entity Role) {
	m.ID = entity.ID
	m.Name = entity.Name
	m.Description = entity.Description
	m.IsActive = entity.IsActive
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}

// PermissionModel represents the database table structure for permissions
type PermissionModel struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string    `gorm:"column:name;uniqueIndex;not null"`
	Resource    string    `gorm:"column:resource;not null"`
	Action      string    `gorm:"column:action;not null"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name for GORM
func (PermissionModel) TableName() string {
	return "permissions"
}

// ToEntity converts a database model to a domain entity
func (m PermissionModel) ToEntity() Permission {
	return Permission{
		ID:          m.ID,
		Name:        m.Name,
		Resource:    m.Resource,
		Action:      m.Action,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *PermissionModel) FromEntity(entity Permission) {
	m.ID = entity.ID
	m.Name = entity.Name
	m.Resource = entity.Resource
	m.Action = entity.Action
	m.Description = entity.Description
	m.CreatedAt = entity.CreatedAt
}

// RolePermissionModel represents the database table structure for role_permissions
type RolePermissionModel struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	RoleID       uuid.UUID `gorm:"column:role_id;not null"`
	PermissionID uuid.UUID `gorm:"column:permission_id;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name for GORM
func (RolePermissionModel) TableName() string {
	return "role_permissions"
}

// ToEntity converts a database model to a domain entity
func (m RolePermissionModel) ToEntity() RolePermission {
	return RolePermission{
		ID:           m.ID,
		RoleID:       m.RoleID,
		PermissionID: m.PermissionID,
		CreatedAt:    m.CreatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *RolePermissionModel) FromEntity(entity RolePermission) {
	m.ID = entity.ID
	m.RoleID = entity.RoleID
	m.PermissionID = entity.PermissionID
	m.CreatedAt = entity.CreatedAt
}