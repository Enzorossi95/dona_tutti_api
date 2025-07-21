package activity

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

// @Summary Get activities for a campaign
// @Description Get all activities for a specific campaign
// @Tags activities
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Success 200 {array} Activity
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{campaignId}/activities [get]
func (h *Handler) GetActivitiesByCampaign(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	activities, err := h.service.GetActivitiesByCampaign(c.Request().Context(), campaignID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, activities)
}

// @Summary Get activity by ID
// @Description Get activity details by ID
// @Tags activities
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Activity ID"
// @Success 200 {object} Activity
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{campaignId}/activities/{id} [get]
func (h *Handler) GetActivity(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid activity ID")
	}

	activity, err := h.service.GetActivity(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, activity)
}

// @Summary Create a new activity
// @Description Create a new activity for a campaign
// @Tags activities
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param activity body Activity true "Activity details"
// @Success 201 {object} Activity
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{campaignId}/activities [post]
func (h *Handler) CreateActivity(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	var activity Activity
	if err := c.Bind(&activity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format: "+err.Error())
	}

	// Set the campaign ID from the URL parameter
	activity.CampaignID = campaignID

	id, err := h.service.CreateActivity(c.Request().Context(), activity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	activity.ID = id
	return c.JSON(http.StatusCreated, activity)
}

// @Summary Update activity details (partial update)
// @Description Update activity details by ID. Only provided fields will be updated.
// @Tags activities
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Activity ID"
// @Param activity body ActivityUpdateRequest true "Activity update details"
// @Success 200 {object} Activity
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{campaignId}/activities/{id} [put]
func (h *Handler) UpdateActivity(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid activity ID")
	}

	var updateReq ActivityUpdateRequest
	if err := c.Bind(&updateReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format: "+err.Error())
	}

	if err := h.service.UpdateActivity(c.Request().Context(), id, campaignID, updateReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get updated activity to return
	updatedActivity, err := h.service.GetActivity(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, updatedActivity)
}

// @Summary Delete an activity
// @Description Delete an activity by ID
// @Tags activities
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Activity ID"
// @Success 204
// @Failure 400 {object} errors.APIError
// @Router /campaigns/{campaignId}/activities/{id} [delete]
func (h *Handler) DeleteActivity(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid activity ID")
	}

	if err := h.service.DeleteActivity(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}