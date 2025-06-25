package campaign

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

// RegisterRoutes registers all campaign routes
func RegisterRoutes(g *echo.Group, service Service) {
	handler := NewHandler(service)

	// Campaign routes
	campaignGroup := g.Group("/campaigns")
	campaignGroup.GET("", handler.ListCampaigns)
	campaignGroup.GET("/:id", handler.GetCampaign)
	campaignGroup.POST("", handler.CreateCampaign)
	campaignGroup.PUT("/:id", handler.UpdateCampaign)
	campaignGroup.DELETE("/:id", handler.DeleteCampaign)
}

// @Summary List all campaigns
// @Description Get a list of all campaigns
// @Tags campaigns
// @Accept json
// @Produce json
// @Success 200 {array} Campaign
// @Failure 400 {object} errors.APIError
// @Router /campaigns [get]
func (h *Handler) ListCampaigns(c echo.Context) error {
	campaigns, err := h.service.ListCampaigns(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, campaigns)
}

// @Summary Get campaign by ID
// @Description Get campaign details by ID
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} Campaign
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{id} [get]
func (h *Handler) GetCampaign(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	campaign, err := h.service.GetCampaign(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, campaign)
}

// @Summary Create a new campaign
// @Description Create a new campaign with the provided details
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign body Campaign true "Campaign details"
// @Success 201 {object} Campaign
// @Failure 400 {object} errors.APIError
// @Router /campaigns [post]
func (h *Handler) CreateCampaign(c echo.Context) error {
	var campaign Campaign
	if err := c.Bind(&campaign); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	id, err := h.service.CreateCampaign(c.Request().Context(), campaign)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	campaign.ID = id
	return c.JSON(http.StatusCreated, campaign)
}

// @Summary Update campaign details
// @Description Update campaign details by ID
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Param campaign body Campaign true "Campaign details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{id} [put]
func (h *Handler) UpdateCampaign(c echo.Context) error {
	// Implementation needed
	return echo.NewHTTPError(http.StatusNotImplemented, "Update campaign details not implemented")
}

// @Summary Delete a campaign
// @Description Delete a campaign by ID
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 204
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{id} [delete]
func (h *Handler) DeleteCampaign(c echo.Context) error {
	// Implementation needed
	return echo.NewHTTPError(http.StatusNotImplemented, "Delete campaign not implemented")
}
