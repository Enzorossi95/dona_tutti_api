package organizer

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

// RegisterRoutes registers all organizer routes
func RegisterRoutes(g *echo.Group, service Service) {
	handler := NewHandler(service)

	// Organizer routes
	organizerGroup := g.Group("/organizers")

	// Public routes (read-only access for everyone)
	organizerGroup.GET("", handler.ListOrganizers)
	organizerGroup.GET("/:id", handler.GetOrganizer)

	// Protected routes requiring authentication and admin role
	authGroup := organizerGroup.Group("", middleware.RequireAuth())

	// Note: RBAC middleware would be initialized here when fully integrated
	// For now, these routes require authentication but not specific roles
	// TODO: Add rbacMiddleware := middleware.NewRBACMiddleware(db.(*gorm.DB))
	// TODO: Add rbacMiddleware.RequireRole("admin") to the routes below

	// Admin-only routes (authentication required, admin role to be added)
	authGroup.POST("", handler.CreateOrganizer)    // Future: Admin only
	authGroup.PUT("/:id", handler.UpdateOrganizer) // Future: Admin only
}

// @Summary List all organizers
// @Description Get a list of all organizers
// @Tags organizers
// @Accept json
// @Produce json
// @Success 200 {array} Organizer
// @Failure 400 {object} errors.APIError
// @Router /organizers [get]
func (h *Handler) ListOrganizers(c echo.Context) error {
	organizers, err := h.service.ListOrganizers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, organizers)
}

// @Summary Get organizer by ID
// @Description Get organizer details by ID
// @Tags organizers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organizer ID"
// @Success 200 {object} Organizer
// @Failure 400 {object} errors.APIError
// @Failure 401 {object} errors.APIError "Unauthorized"
// @Router /organizers/{id} [get]
func (h *Handler) GetOrganizer(c echo.Context) error {
	// Opcionalmente, podemos obtener el ID del usuario del token
	//userID := c.Get("user_id").(string)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organizer ID")
	}

	organizer, err := h.service.GetOrganizer(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, organizer)
}

// @Summary Create a new organizer
// @Description Create a new organizer with the provided details
// @Tags organizers
// @Accept json
// @Produce json
// @Param organizer body Organizer true "Organizer details"
// @Success 201 {object} Organizer
// @Failure 400 {object} errors.APIError
// @Router /organizers [post]
func (h *Handler) CreateOrganizer(c echo.Context) error {
	// Implementation of CreateOrganizer method
	return nil // Placeholder return, actual implementation needed
}

// @Summary Update organizer details
// @Description Update organizer details by ID
// @Tags organizers
// @Accept json
// @Produce json
// @Param id path string true "Organizer ID"
// @Param organizer body Organizer true "Organizer details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /organizers/{id} [put]
func (h *Handler) UpdateOrganizer(c echo.Context) error {
	// Implementation of UpdateOrganizer method
	return nil // Placeholder return, actual implementation needed
}
