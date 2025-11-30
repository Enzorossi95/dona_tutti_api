package middleware

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AuthContext represents the authentication context for a user
type AuthContext struct {
	UserID      uuid.UUID `json:"user_id"`
	RoleID      uuid.UUID `json:"role_id"`
	RoleName    string    `json:"role_name"`
	Permissions []string  `json:"permissions"`
}

// RBACService defines the interface for RBAC operations needed by middleware
type RBACService interface {
	// Role operations
	HasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error)
	HasAnyRole(ctx context.Context, userID uuid.UUID, roleNames []string) (bool, error)

	// Permission operations
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)

	// Ownership operations
	ValidateResourceOwnership(ctx context.Context, userID uuid.UUID, resource string, resourceID uuid.UUID) (bool, error)

	// User context operations
	GetUserAuthContext(ctx context.Context, userID uuid.UUID) (interface{}, error) // Returns implementation-specific auth context
}

// RBACMiddleware provides role-based access control functionality
type RBACMiddleware struct {
	rbacService RBACService
}

// NewRBACMiddleware creates a new RBAC middleware instance
func NewRBACMiddleware(rbacService RBACService) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
	}
}

// RoleConfig represents configuration for role-based authorization
type RoleConfig struct {
	Roles          []string
	RequireAll     bool // If true, user must have ALL roles; if false, user must have ANY role
	AllowOwnership bool // If true, resource ownership can bypass role requirements
}

// PermissionConfig represents configuration for permission-based authorization
type PermissionConfig struct {
	Permissions    []string
	RequireAll     bool // If true, user must have ALL permissions; if false, user must have ANY permission
	AllowOwnership bool // If true, resource ownership can bypass permission requirements
}

// OwnershipConfig represents configuration for ownership-based authorization
type OwnershipConfig struct {
	Resource         string
	ResourceIDParam  string // URL parameter name for resource ID (e.g., "id", "user_id")
	UserIDField      string // Field name in resource that contains user ID
	AllowAdminBypass bool   // If true, admin role can bypass ownership check
}

// RequireRole creates middleware that requires specific role(s)
func (m *RBACMiddleware) RequireRole(roles ...string) echo.MiddlewareFunc {
	return m.RequireRoleWithConfig(RoleConfig{
		Roles:      roles,
		RequireAll: false,
	})
}

// RequireAllRoles creates middleware that requires ALL specified roles
func (m *RBACMiddleware) RequireAllRoles(roles ...string) echo.MiddlewareFunc {
	return m.RequireRoleWithConfig(RoleConfig{
		Roles:      roles,
		RequireAll: true,
	})
}

// RequireRoleWithConfig creates middleware with custom role configuration
func (m *RBACMiddleware) RequireRoleWithConfig(config RoleConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := m.getUserIDFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			// Check role requirements
			if config.RequireAll {
				for _, role := range config.Roles {
					hasRole, err := m.rbacService.HasRole(c.Request().Context(), userID, role)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, "Authorization check failed")
					}
					if !hasRole {
						if config.AllowOwnership {
							if allowed, _ := m.checkOwnership(c, userID); allowed {
								break
							}
						}
						return echo.NewHTTPError(http.StatusForbidden, "Insufficient role privileges")
					}
				}
			} else {
				hasAnyRole, err := m.rbacService.HasAnyRole(c.Request().Context(), userID, config.Roles)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Authorization check failed")
				}
				if !hasAnyRole {
					if config.AllowOwnership {
						if allowed, _ := m.checkOwnership(c, userID); !allowed {
							return echo.NewHTTPError(http.StatusForbidden, "Insufficient role privileges")
						}
					} else {
						return echo.NewHTTPError(http.StatusForbidden, "Insufficient role privileges")
					}
				}
			}

			return next(c)
		}
	}
}

// RequirePermission creates middleware that requires specific permission(s)
func (m *RBACMiddleware) RequirePermission(permissions ...string) echo.MiddlewareFunc {
	return m.RequirePermissionWithConfig(PermissionConfig{
		Permissions: permissions,
		RequireAll:  false,
	})
}

// RequireAllPermissions creates middleware that requires ALL specified permissions
func (m *RBACMiddleware) RequireAllPermissions(permissions ...string) echo.MiddlewareFunc {
	return m.RequirePermissionWithConfig(PermissionConfig{
		Permissions: permissions,
		RequireAll:  true,
	})
}

// RequirePermissionWithConfig creates middleware with custom permission configuration
func (m *RBACMiddleware) RequirePermissionWithConfig(config PermissionConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := m.getUserIDFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			// Check permission requirements
			if config.RequireAll {
				for _, permission := range config.Permissions {
					hasPermission, err := m.rbacService.HasPermission(c.Request().Context(), userID, permission)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, "Authorization check failed")
					}
					if !hasPermission {
						if config.AllowOwnership {
							if allowed, _ := m.checkOwnership(c, userID); allowed {
								break
							}
						}
						return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
					}
				}
			} else {
				hasAnyPermission := false
				for _, permission := range config.Permissions {
					hasPermission, err := m.rbacService.HasPermission(c.Request().Context(), userID, permission)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, "Authorization check failed")
					}
					if hasPermission {
						hasAnyPermission = true
						break
					}
				}

				if !hasAnyPermission {
					if config.AllowOwnership {
						if allowed, _ := m.checkOwnership(c, userID); !allowed {
							return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
						}
					} else {
						return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
					}
				}
			}

			return next(c)
		}
	}
}

// RequireOwnership creates middleware that requires resource ownership
func (m *RBACMiddleware) RequireOwnership() echo.MiddlewareFunc {
	return m.RequireOwnershipWithConfig(OwnershipConfig{
		ResourceIDParam:  "id",
		AllowAdminBypass: true,
	})
}

// RequireOwnershipWithConfig creates middleware with custom ownership configuration
func (m *RBACMiddleware) RequireOwnershipWithConfig(config OwnershipConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := m.getUserIDFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			// Check if admin bypass is allowed
			if config.AllowAdminBypass {
				isAdmin, err := m.rbacService.HasRole(c.Request().Context(), userID, "admin")
				if err == nil && isAdmin {
					return next(c)
				}
			}

			// Check ownership
			allowed, err := m.checkOwnershipWithConfig(c, userID, config)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Ownership check failed")
			}
			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied: resource ownership required")
			}

			return next(c)
		}
	}
}

// RequireAdminRole creates middleware that requires admin role
func (m *RBACMiddleware) RequireAdminRole() echo.MiddlewareFunc {
	return m.RequireRole("admin")
}

// RequireAnyRole creates middleware that allows any of the specified roles
func (m *RBACMiddleware) RequireAnyRole(roles ...string) echo.MiddlewareFunc {
	return m.RequireRoleWithConfig(RoleConfig{
		Roles:      roles,
		RequireAll: false,
	})
}

// Combine creates middleware that combines multiple authorization checks with OR logic
func (m *RBACMiddleware) Combine(middlewares ...echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try each middleware in sequence; if any passes, allow access
			for _, middleware := range middlewares {
				// Create a test context with the same request and response
				testContext := c.Echo().NewContext(c.Request(), c.Response())

				// Copy all context values from the original context
				if userID := c.Get("user_id"); userID != nil {
					testContext.Set("user_id", userID)
				}
				if roleID := c.Get("role_id"); roleID != nil {
					testContext.Set("role_id", roleID)
				}
				if role := c.Get("role"); role != nil {
					testContext.Set("role", role)
				}

				err := middleware(func(c echo.Context) error { return nil })(testContext)
				if err == nil {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}
	}
}

// Helper functions

func (m *RBACMiddleware) getUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	// Handle both string and uuid.UUID types
	switch v := userIDValue.(type) {
	case string:
		userID, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID format")
		}
		return userID, nil
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID type")
	}
}

func (m *RBACMiddleware) checkOwnership(c echo.Context, userID uuid.UUID) (bool, error) {
	// Default ownership check - user can access their own user resource
	resourceID := c.Param("id")
	if resourceID == "" {
		return false, nil
	}

	resourceUUID, err := uuid.Parse(resourceID)
	if err != nil {
		return false, err
	}

	// For user resources, check if the user is accessing their own resource
	if strings.Contains(c.Path(), "/users/") {
		return userID == resourceUUID, nil
	}

	return false, nil
}

func (m *RBACMiddleware) checkOwnershipWithConfig(c echo.Context, userID uuid.UUID, config OwnershipConfig) (bool, error) {
	resourceID := c.Param(config.ResourceIDParam)
	if resourceID == "" {
		return false, nil
	}

	resourceUUID, err := uuid.Parse(resourceID)
	if err != nil {
		return false, err
	}

	if config.Resource != "" {
		return m.rbacService.ValidateResourceOwnership(c.Request().Context(), userID, config.Resource, resourceUUID)
	}

	// Default behavior - check if accessing own user resource
	return userID == resourceUUID, nil
}

// RequireRoleMiddleware is a helper function that creates a middleware requiring a specific role
// This is useful when you need to create middleware without first creating an RBACMiddleware instance
func RequireRoleMiddleware(rbacService RBACService, role string) echo.MiddlewareFunc {
	m := NewRBACMiddleware(rbacService)
	return m.RequireRole(role)
}

// SetUserContext sets the user context in the Echo context for RBAC operations
func (m *RBACMiddleware) SetUserContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := m.getUserIDFromContext(c)
			if err != nil {
				return next(c) // Continue without user context
			}

			// Get user auth context and store it
			authCtxInterface, err := m.rbacService.GetUserAuthContext(c.Request().Context(), userID)
			if err == nil && authCtxInterface != nil {
				// Use reflection to extract fields from the returned struct
				authCtx := &AuthContext{}
				v := reflect.ValueOf(authCtxInterface)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				
				if v.Kind() == reflect.Struct {
					// Extract UserID field
					if userIDField := v.FieldByName("UserID"); userIDField.IsValid() && userIDField.CanInterface() {
						if uid, ok := userIDField.Interface().(uuid.UUID); ok {
							authCtx.UserID = uid
						}
					}
					
					// Extract RoleID field
					if roleIDField := v.FieldByName("RoleID"); roleIDField.IsValid() && roleIDField.CanInterface() {
						if rid, ok := roleIDField.Interface().(uuid.UUID); ok {
							authCtx.RoleID = rid
						}
					}
					
					// Extract RoleName field
					if roleNameField := v.FieldByName("RoleName"); roleNameField.IsValid() && roleNameField.CanInterface() {
						if rn, ok := roleNameField.Interface().(string); ok {
							authCtx.RoleName = rn
						}
					}
					
					// Extract Permissions field
					if permissionsField := v.FieldByName("Permissions"); permissionsField.IsValid() && permissionsField.CanInterface() {
						if perms, ok := permissionsField.Interface().([]string); ok {
							authCtx.Permissions = perms
						}
					}
					
					c.Set("auth_context", authCtx)
					c.Set("user_role", authCtx.RoleName)
					c.Set("user_permissions", authCtx.Permissions)
				}
			}

			return next(c)
		}
	}
}
