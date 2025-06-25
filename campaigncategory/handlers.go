package campaigncategory

import (
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

// RegisterRoutes registers all category routes
func RegisterRoutes(g *echo.Group, service Service) {
	handler := NewHandler(service)

	// Category routes
	categoryGroup := g.Group("/categories")
	categoryGroup.GET("", handler.ListCategories)
	categoryGroup.GET("/:id", handler.GetCategory)
	categoryGroup.POST("", handler.CreateCategory)
	categoryGroup.PUT("/:id", handler.UpdateCategory)
}

// @Summary List all categories
// @Description Get a list of all campaign categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} CampaignCategory
// @Failure 400 {object} errors.APIError
// @Router /categories [get]
func (h *Handler) ListCategories(c echo.Context) error {
	categories, err := h.service.ListCategories(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}

// @Summary Get category by ID
// @Description Get category details by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} CampaignCategory
// @Failure 400 {object} errors.APIError
// @Router /categories/{id} [get]
func (h *Handler) GetCategory(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid category ID")
	}

	category, err := h.service.GetCategory(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, category)
}

// @Summary Create a new category
// @Description Create a new campaign category with the provided details
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CampaignCategory true "Category details"
// @Success 201 {object} CampaignCategory
// @Failure 400 {object} errors.APIError
// @Router /categories [post]
func (h *Handler) CreateCategory(c echo.Context) error {
	// Implementation of CreateCategory method
	return nil // Placeholder return, actual implementation needed
}

// @Summary Update category details
// @Description Update category details by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body CampaignCategory true "Category details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(c echo.Context) error {
	// Implementation of UpdateCategory method
	return nil // Placeholder return, actual implementation needed
}
