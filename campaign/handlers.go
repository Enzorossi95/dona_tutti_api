package campaign

import (
	"dona_tutti_api/campaign/activity"
	"dona_tutti_api/middleware"
	"dona_tutti_api/s3client"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  Service
	s3Client *s3client.Client
}

func NewHandler(service Service, s3Client *s3client.Client) *Handler {
	return &Handler{service: service, s3Client: s3Client}
}

// RegisterRoutes registers all campaign routes with RBAC authorization
func RegisterRoutes(g *echo.Group, service Service, activityService activity.Service, s3Client *s3client.Client, rbacService middleware.RBACService) {
	handler := NewHandler(service, s3Client)
	activityHandler := activity.NewHandler(activityService)
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	// Campaign routes
	campaignGroup := g.Group("/campaigns")

	// Public routes (no authentication required)
	campaignGroup.GET("/:id", handler.GetCampaign)
	campaignGroup.GET("/:campaignId/activities", activityHandler.GetActivitiesByCampaign)
	campaignGroup.GET("/:campaignId/activities/:id", activityHandler.GetActivity)

	// Protected routes with authentication
	authGroup := campaignGroup.Group("", middleware.RequireAuth())

	// Admin-only routes
	adminGroup := authGroup.Group("", rbacMiddleware.RequireRole("admin"))
	adminGroup.POST("", handler.CreateCampaign)
	adminGroup.DELETE("/:id", handler.DeleteCampaign)
	adminGroup.GET("", handler.ListCampaigns)

	// Campaign upload routes
	adminGroup.POST("/:id/upload", handler.UploadCampaignImage)

	// Activity admin routes
	adminGroup.POST("/:campaignId/activities", activityHandler.CreateActivity)
	adminGroup.PUT("/:campaignId/activities/:id", activityHandler.UpdateActivity)
	adminGroup.DELETE("/:campaignId/activities/:id", activityHandler.DeleteActivity)

	// Admin or owner routes (using Combine for OR logic)
	adminOrOwnerGroup := authGroup.Group("", rbacMiddleware.Combine(
		rbacMiddleware.RequireRole("admin"),
		rbacMiddleware.RequireOwnership(),
	))
	adminOrOwnerGroup.PUT("/:id", handler.UpdateCampaign)
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format: "+err.Error())
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

// @Summary Upload campaign image
// @Description Upload an image for a campaign
// @Tags campaigns
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Campaign ID"
// @Param file formData file true "Image file"
// @Success 200 {object} s3client.UploadResponse
// @Failure 400 {object} errors.APIError
// @Failure 503 {object} errors.APIError
// @Router /campaigns/{id}/upload [post]
func (h *Handler) UploadCampaignImage(c echo.Context) error {
	// Check if S3 client is available
	if h.s3Client == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, 
			"File upload service not available. AWS S3 not configured.")
	}

	campaignID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No file provided")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to open file")
	}
	defer src.Close()

	// Create upload request
	uploadReq := s3client.UploadRequest{
		File:         src,
		Header:       file,
		ResourceType: "campaign",
		ResourceID:   campaignID.String(),
	}

	// Upload to S3
	response, err := h.s3Client.Upload(c.Request().Context(), uploadReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update campaign image URL in database
	if err := h.service.UpdateCampaignImage(c.Request().Context(), campaignID, response.URL); err != nil {
		// If database update fails, try to delete the uploaded file
		h.s3Client.DeleteByKey(c.Request().Context(), response.Key)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
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
