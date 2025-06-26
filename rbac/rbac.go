package rbac

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role in the system
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission represents a specific permission in the system
type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           uuid.UUID `json:"id"`
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserRole represents a user with their role information
type UserRole struct {
	UserID   uuid.UUID `json:"user_id"`
	RoleID   uuid.UUID `json:"role_id"`
	RoleName string    `json:"role_name"`
}

// PermissionCheck represents a permission check request
type PermissionCheck struct {
	UserID     uuid.UUID `json:"user_id"`
	Resource   string    `json:"resource"`
	Action     string    `json:"action"`
	ResourceID *string   `json:"resource_id,omitempty"` // For contextual permissions
}

// Default role constants
const (
	RoleAdmin = "admin"
	RoleDonor = "donor"
	RoleGuest = "guest"
)

// Default role UUIDs (matching migration)
var (
	AdminRoleID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	DonorRoleID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	GuestRoleID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
)

// Permission constants
const (
	// Campaign permissions
	PermissionCampaignsCreate = "campaigns:create"
	PermissionCampaignsRead   = "campaigns:read"
	PermissionCampaignsUpdate = "campaigns:update"
	PermissionCampaignsDelete = "campaigns:delete"

	// Donation permissions
	PermissionDonationsCreate = "donations:create"
	PermissionDonationsRead   = "donations:read"
	PermissionDonationsUpdate = "donations:update"
	PermissionDonationsDelete = "donations:delete"

	// User permissions
	PermissionUsersCreate = "users:create"
	PermissionUsersRead   = "users:read"
	PermissionUsersUpdate = "users:update"
	PermissionUsersDelete = "users:delete"

	// Category permissions
	PermissionCategoriesCreate = "categories:create"
	PermissionCategoriesRead   = "categories:read"
	PermissionCategoriesUpdate = "categories:update"
	PermissionCategoriesDelete = "categories:delete"

	// Organizer permissions
	PermissionOrganizersCreate = "organizers:create"
	PermissionOrganizersRead   = "organizers:read"
	PermissionOrganizersUpdate = "organizers:update"
	PermissionOrganizersDelete = "organizers:delete"

	// Donor permissions
	PermissionDonorsCreate = "donors:create"
	PermissionDonorsRead   = "donors:read"
	PermissionDonorsUpdate = "donors:update"
	PermissionDonorsDelete = "donors:delete"
)

// Resource constants
const (
	ResourceCampaigns = "campaigns"
	ResourceDonations = "donations"
	ResourceUsers     = "users"
	ResourceCategories = "categories"
	ResourceOrganizers = "organizers"
	ResourceDonors    = "donors"
)

// Action constants
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
)