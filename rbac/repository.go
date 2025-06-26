package rbac

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocalAuthContext represents the authentication context for a user (local to rbac package)
type LocalAuthContext struct {
	UserID      uuid.UUID `json:"user_id"`
	RoleID      uuid.UUID `json:"role_id"`
	RoleName    string    `json:"role_name"`
	Permissions []string  `json:"permissions"`
}

type Repository interface {
	// Role operations
	GetRoleByID(ctx context.Context, id uuid.UUID) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
	CreateRole(ctx context.Context, role Role) error
	UpdateRole(ctx context.Context, role Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error

	// Permission operations
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*Permission, error)
	ListPermissions(ctx context.Context) ([]Permission, error)
	HasPermission(ctx context.Context, roleID uuid.UUID, permissionName string) (bool, error)

	// Role-Permission operations
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error

	// User context operations
	GetUserAuthContext(ctx context.Context, userID uuid.UUID) (interface{}, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Role operations
func (r *repository) GetRoleByID(ctx context.Context, id uuid.UUID) (*Role, error) {
	var model RoleModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	role := model.ToEntity()
	return &role, nil
}

func (r *repository) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	var model RoleModel
	if err := r.db.WithContext(ctx).Where("name = ? AND is_active = ?", name, true).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	role := model.ToEntity()
	return &role, nil
}

func (r *repository) ListRoles(ctx context.Context) ([]Role, error) {
	var models []RoleModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	roles := make([]Role, len(models))
	for i, model := range models {
		roles[i] = model.ToEntity()
	}
	return roles, nil
}

func (r *repository) CreateRole(ctx context.Context, role Role) error {
	var model RoleModel
	model.FromEntity(role)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

func (r *repository) UpdateRole(ctx context.Context, role Role) error {
	var model RoleModel
	model.FromEntity(role)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (r *repository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&RoleModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// Permission operations
func (r *repository) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]Permission, error) {
	var permissions []Permission
	
	query := `
		SELECT p.id, p.name, p.resource, p.action, p.description, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`
	
	if err := r.db.WithContext(ctx).Raw(query, roleID).Scan(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get permissions for role: %w", err)
	}
	
	return permissions, nil
}

func (r *repository) GetPermissionByName(ctx context.Context, name string) (*Permission, error) {
	var model PermissionModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	permission := model.ToEntity()
	return &permission, nil
}

func (r *repository) ListPermissions(ctx context.Context) ([]Permission, error) {
	var models []PermissionModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	permissions := make([]Permission, len(models))
	for i, model := range models {
		permissions[i] = model.ToEntity()
	}
	return permissions, nil
}

func (r *repository) HasPermission(ctx context.Context, roleID uuid.UUID, permissionName string) (bool, error) {
	var count int64
	
	query := `
		SELECT COUNT(*)
		FROM role_permissions rp
		INNER JOIN permissions p ON rp.permission_id = p.id
		WHERE rp.role_id = ? AND p.name = ?
	`
	
	if err := r.db.WithContext(ctx).Raw(query, roleID, permissionName).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	
	return count > 0, nil
}

// Role-Permission operations
func (r *repository) AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	rolePermission := RolePermissionModel{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	
	if err := r.db.WithContext(ctx).Create(&rolePermission).Error; err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}
	return nil
}

func (r *repository) RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermissionModel{}).Error; err != nil {
		return fmt.Errorf("failed to revoke permission from role: %w", err)
	}
	return nil
}

// User context operations
func (r *repository) GetUserAuthContext(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	var result struct {
		UserID   uuid.UUID `gorm:"column:user_id"`
		RoleID   uuid.UUID `gorm:"column:role_id"`
		RoleName string    `gorm:"column:role_name"`
	}
	
	query := `
		SELECT u.id as user_id, u.role_id, r.name as role_name
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.id = ? AND r.is_active = true
	`
	
	if err := r.db.WithContext(ctx).Raw(query, userID).Scan(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user auth context: %w", err)
	}
	
	// Get user permissions
	permissions, err := r.GetPermissionsByRoleID(ctx, result.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	
	permissionNames := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionNames[i] = permission.Name
	}
	
	return &LocalAuthContext{
		UserID:      result.UserID,
		RoleID:      result.RoleID,
		RoleName:    result.RoleName,
		Permissions: permissionNames,
	}, nil
}