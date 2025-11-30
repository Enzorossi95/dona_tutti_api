package closure

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for campaign closure
type Handler struct {
	service Service
}

// NewHandler creates a new closure handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the closure routes
func (h *Handler) RegisterRoutes(g *echo.Group, authMiddleware echo.MiddlewareFunc, adminMiddleware echo.MiddlewareFunc) {
	// Public routes (no authentication required)
	g.GET("/campaigns/:id/audit", h.GetPublicAuditReport)
	g.GET("/campaigns/:id/audit/download", h.DownloadAuditPDF)

	// Protected routes - require authentication first, then admin role
	authGroup := g.Group("", authMiddleware)
	adminGroup := authGroup.Group("", adminMiddleware)
	adminGroup.POST("/campaigns/:id/close", h.CloseCampaign)
	adminGroup.GET("/campaigns/:id/closure-report", h.GetClosureReport)
}

// CloseCampaignRequestDTO represents the request to close a campaign
type CloseCampaignRequestDTO struct {
	ClosureType string `json:"closure_type" validate:"required,oneof=goal_reached end_date manual"`
	Reason      string `json:"reason"`
}

// CloseCampaign handles POST /api/campaigns/:id/close
// @Summary Close a campaign and generate audit report
// @Description Closes a campaign, generates transparency score and audit report PDF
// @Tags campaign-closure
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Param request body CloseCampaignRequestDTO true "Closure request"
// @Success 200 {object} CampaignClosureReport "Closure report"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Campaign not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/campaigns/{id}/close [post]
func (h *Handler) CloseCampaign(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Parse request body
	var reqDTO CloseCampaignRequestDTO
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// Validate closure type
	closureType := ClosureType(reqDTO.ClosureType)
	if closureType != ClosureTypeGoalReached && closureType != ClosureTypeEndDate && closureType != ClosureTypeManual {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid closure type. Must be: goal_reached, end_date, or manual",
		})
	}

	// Validate reason for manual closure
	if closureType == ClosureTypeManual && len(reqDTO.Reason) < 10 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Manual closure requires a reason with at least 10 characters",
		})
	}

	// Get admin user ID from context (assuming it's set by auth middleware)
	var closedBy *uuid.UUID
	if userID, ok := c.Get("user_id").(uuid.UUID); ok {
		closedBy = &userID
	}

	// Set reason pointer
	var reason *string
	if reqDTO.Reason != "" {
		reason = &reqDTO.Reason
	}

	// Close campaign
	report, err := h.service.CloseCampaign(c.Request().Context(), campaignID, closureType, reason, closedBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, report)
}

// GetClosureReport handles GET /api/campaigns/:id/closure-report
// @Summary Get closure report for a campaign (admin)
// @Description Retrieves the full closure report with all metrics
// @Tags campaign-closure
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} CampaignClosureReport "Closure report"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Report not found"
// @Router /api/campaigns/{id}/closure-report [get]
func (h *Handler) GetClosureReport(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Get closure report
	report, err := h.service.GetClosureReport(c.Request().Context(), campaignID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Closure report not found",
		})
	}

	return c.JSON(http.StatusOK, report)
}

// GetPublicAuditReport handles GET /api/campaigns/:id/audit
// @Summary Get public audit report for donors
// @Description Retrieves the public audit report visible to donors
// @Tags campaign-closure
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} PublicAuditReport "Public audit report"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Report not found"
// @Router /api/campaigns/{id}/audit [get]
func (h *Handler) GetPublicAuditReport(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Get public audit report
	report, err := h.service.GetPublicAuditReport(c.Request().Context(), campaignID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Audit report not found for this campaign",
		})
	}

	return c.JSON(http.StatusOK, report)
}

// DownloadAuditPDF handles GET /api/campaigns/:id/audit/download
// @Summary Download audit report PDF
// @Description Redirects to the PDF URL for download
// @Tags campaign-closure
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 302 "Redirect to PDF URL"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Report not found"
// @Router /api/campaigns/{id}/audit/download [get]
func (h *Handler) DownloadAuditPDF(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Get closure report to get PDF URL
	report, err := h.service.GetClosureReport(c.Request().Context(), campaignID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Audit report not found",
		})
	}

	// Check if PDF is available
	if report.ReportPdfURL == nil || *report.ReportPdfURL == "" {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "PDF report is still being generated",
			"message": "Please try again in a few moments",
		})
	}

	// Redirect to PDF URL
	return c.Redirect(http.StatusFound, *report.ReportPdfURL)
}
