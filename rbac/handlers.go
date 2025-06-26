package rbac

import (
	"dona_tutti_api/middleware"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all RBAC management routes
func RegisterRoutes(g *echo.Group, service Service) {
	handler := NewHandler(service)

	// TEST ROUTE - Simple route without middleware
	g.GET("/rbac-test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "RBAC routes are working!"})
	})

	rbacMiddleware := middleware.NewRBACMiddleware(service)

	// Role management routes (admin only)
	roleGroup := g.Group("/roles", middleware.RequireAuth(), rbacMiddleware.RequireRole("admin"))
	roleGroup.GET("", handler.ListRoles)
	roleGroup.GET("/:name", handler.GetRoleByName)

	// Permission management routes (admin only)
	permissionGroup := g.Group("/permissions", middleware.RequireAuth(), rbacMiddleware.RequireRole("admin"))
	permissionGroup.GET("", handler.ListPermissions)

	// User permission check routes (require authentication)
	authGroup := g.Group("/auth", middleware.RequireAuth(), rbacMiddleware.RequireRole("admin"))
	authGroup.GET("/check-permission", handler.CheckPermission)
	authGroup.GET("/user-context", handler.GetUserContext)
}

// @Summary List all roles
// @Description Get a list of all active roles in the system
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Role
// @Failure 401 {object} errors.APIError
// @Failure 403 {object} errors.APIError
// @Router /roles [get]
func (h *Handler) ListRoles(c echo.Context) error {
	roles, err := h.service.ListRoles(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve roles")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": roles,
	})
}

// @Summary Get role by name
// @Description Get role details by role name
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name path string true "Role name"
// @Success 200 {object} Role
// @Failure 401 {object} errors.APIError
// @Failure 403 {object} errors.APIError
// @Failure 404 {object} errors.APIError
// @Router /roles/{name} [get]
func (h *Handler) GetRoleByName(c echo.Context) error {
	roleName := c.Param("name")
	if roleName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Role name is required")
	}

	role, err := h.service.GetRoleByName(c.Request().Context(), roleName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Role not found")
	}

	return c.JSON(http.StatusOK, role)
}

// @Summary List all permissions
// @Description Get a list of all permissions in the system
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Permission
// @Failure 401 {object} errors.APIError
// @Failure 403 {object} errors.APIError
// @Router /permissions [get]
func (h *Handler) ListPermissions(c echo.Context) error {
	permissions, err := h.service.ListPermissions(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve permissions")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"permissions": permissions,
	})
}

// CheckPermissionRequest represents a permission check request
type CheckPermissionRequest struct {
	Permission string `json:"permission" validate:"required"`
}

// @Summary Check user permission
// @Description Check if the current user has a specific permission
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CheckPermissionRequest true "Permission check request"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} errors.APIError
// @Failure 401 {object} errors.APIError
// @Router /auth/check-permission [get]
func (h *Handler) CheckPermission(c echo.Context) error {
	// Get user ID from context (set by auth middleware)
	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
	}

	// Get permission from query parameter
	permission := c.QueryParam("permission")
	if permission == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Permission parameter is required")
	}

	hasPermission, err := h.service.HasPermission(c.Request().Context(), userID, permission)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check permission")
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"has_permission": hasPermission,
	})
}

// @Summary Get user context
// @Description Get the current user's authentication context including role and permissions
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} LocalAuthContext
// @Failure 401 {object} errors.APIError
// @Router /auth/user-context [get]
func (h *Handler) GetUserContext(c echo.Context) error {
	// Get user ID from context (set by auth middleware)
	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
	}

	authContext, err := h.service.GetUserAuthContext(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user context")
	}

	return c.JSON(http.StatusOK, authContext)
}
