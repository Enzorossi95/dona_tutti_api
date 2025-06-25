package donor

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

// RegisterRoutes registers all donor routes
func RegisterRoutes(g *echo.Group, service Service) {
	handler := NewHandler(service)

	// Donor routes
	donorGroup := g.Group("/donors")
	donorGroup.GET("", handler.ListDonors)
	donorGroup.GET("/:id", handler.GetDonor)
	donorGroup.POST("", handler.CreateDonor)
	donorGroup.PUT("/:id", handler.UpdateDonor)
}

// @Summary List all donors
// @Description Get a list of all donors
// @Tags donors
// @Accept json
// @Produce json
// @Success 200 {array} Donor
// @Failure 400 {object} errors.APIError
// @Router /donors [get]
func (h *Handler) ListDonors(c echo.Context) error {
	donors, err := h.service.ListDonors(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, donors)
}

// @Summary Get donor by ID
// @Description Get donor details by ID
// @Tags donors
// @Accept json
// @Produce json
// @Param id path string true "Donor ID"
// @Success 200 {object} Donor
// @Failure 400 {object} errors.APIError
// @Router /donors/{id} [get]
func (h *Handler) GetDonor(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid donor ID")
	}

	donor, err := h.service.GetDonor(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, donor)
}

// @Summary Create a new donor
// @Description Create a new donor with the provided details
// @Tags donors
// @Accept json
// @Produce json
// @Param donor body Donor true "Donor details"
// @Success 201 {object} Donor
// @Failure 400 {object} errors.APIError
// @Router /donors [post]
func (h *Handler) CreateDonor(c echo.Context) error {
	var donor Donor
	if err := c.Bind(&donor); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	id, err := h.service.CreateDonor(c.Request().Context(), donor)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	donor.ID = id
	return c.JSON(http.StatusCreated, donor)
}

// @Summary Update donor details
// @Description Update donor details by ID
// @Tags donors
// @Accept json
// @Produce json
// @Param id path string true "Donor ID"
// @Param donor body Donor true "Donor details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /donors/{id} [put]
func (h *Handler) UpdateDonor(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid donor ID")
	}

	var donor Donor
	if err := c.Bind(&donor); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	donor.ID = id
	err = h.service.UpdateDonor(c.Request().Context(), donor)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
