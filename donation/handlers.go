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

// GetDonationsByCampaign returns all donations for a campaign
// @Summary Get all donations for a campaign
// @Description Get a list of all donations for a specific campaign
// @Tags donations
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Success 200 {array} Donation
// @Failure 400 {object} map[string]string
// @Router /campaigns/{campaignId}/donations [get]
func (h *Handler) GetDonationsByCampaign(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	donations, err := h.service.ListDonationsByCampaign(c.Request().Context(), campaignID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, donations)
}

// GetDonation returns a specific donation
// @Summary Get a donation by ID
// @Description Get donation details by ID
// @Tags donations
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Donation ID"
// @Success 200 {object} Donation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /campaigns/{campaignId}/donations/{id} [get]
func (h *Handler) GetDonation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid donation ID")
	}

	donation, err := h.service.GetDonation(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, donation)
}

// CreateDonation creates a new donation
// @Summary Create a new donation
// @Description Create a new donation for a campaign
// @Tags donations
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param donation body Donation true "Donation data"
// @Success 201 {object} Donation
// @Failure 400 {object} map[string]string
// @Router /campaigns/{campaignId}/donations [post]
func (h *Handler) CreateDonation(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	var donation Donation
	if err := c.Bind(&donation); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Ensure the donation is for the correct campaign
	donation.CampaignID = campaignID

	id, err := h.service.CreateDonation(c.Request().Context(), donation)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	donation.ID = id
	return c.JSON(http.StatusCreated, donation)
}

// UpdateDonation updates an existing donation
// @Summary Update a donation
// @Description Update an existing donation
// @Tags donations
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Donation ID"
// @Param donation body Donation true "Donation update data"
// @Success 200 {object} Donation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /campaigns/{campaignId}/donations/{id} [put]
func (h *Handler) UpdateDonation(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid donation ID")
	}

	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	var donation Donation
	if err := c.Bind(&donation); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Ensure the donation IDs are correct
	donation.ID = id
	donation.CampaignID = campaignID

	if err := h.service.UpdateDonation(c.Request().Context(), donation); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, donation)
}