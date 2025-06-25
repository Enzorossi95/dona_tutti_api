package donation

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

// RegisterRoutes registers all donation routes
func RegisterRoutes(g *echo.Group, service Service) {
	handler := NewHandler(service)

	// Donation routes
	donationGroup := g.Group("/donations")
	donationGroup.GET("", handler.ListDonations)
	donationGroup.GET("/:id", handler.GetDonation)
	donationGroup.POST("", handler.CreateDonation)
	donationGroup.PUT("/:id", handler.UpdateDonation)
}

// @Summary List all donations
// @Description Get a list of all donations
// @Tags donations
// @Accept json
// @Produce json
// @Success 200 {array} Donation
// @Failure 400 {object} errors.APIError
// @Router /donations [get]
func (h *Handler) ListDonations(c echo.Context) error {
	donations, err := h.service.ListDonations(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, donations)
}

// @Summary Get donation by ID
// @Description Get donation details by ID
// @Tags donations
// @Accept json
// @Produce json
// @Param id path string true "Donation ID"
// @Success 200 {object} Donation
// @Failure 400 {object} errors.APIError
// @Router /donations/{id} [get]
func (h *Handler) GetDonation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid donation ID")
	}

	donation, err := h.service.GetDonation(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, donation)
}

// @Summary Create a new donation
// @Description Create a new donation with the provided details
// @Tags donations
// @Accept json
// @Produce json
// @Param donation body Donation true "Donation details"
// @Success 201 {object} Donation
// @Failure 400 {object} errors.APIError
// @Router /donations [post]
func (h *Handler) CreateDonation(c echo.Context) error {
	var donation Donation
	if err := c.Bind(&donation); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	id, err := h.service.CreateDonation(c.Request().Context(), donation)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	donation.ID = id
	return c.JSON(http.StatusCreated, donation)
}

// @Summary Update donation details
// @Description Update donation details by ID
// @Tags donations
// @Accept json
// @Produce json
// @Param id path string true "Donation ID"
// @Param donation body Donation true "Donation details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /donations/{id} [put]
func (h *Handler) UpdateDonation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid donation ID")
	}

	var donation Donation
	if err := c.Bind(&donation); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	donation.ID = id
	err = h.service.UpdateDonation(c.Request().Context(), donation)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// @Summary List donations by campaign
// @Description Get a list of donations for a specific campaign
// @Tags donations
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {array} Donation
// @Failure 400 {object} errors.APIError
// @Router /api/donations/campaign/{id} [get]
func (h *Handler) ListDonationsByCampaign(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	donations, err := h.service.ListDonationsByCampaign(c.Request().Context(), campaignID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, donations)
}
