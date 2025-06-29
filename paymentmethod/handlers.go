package paymentmethod

import (
	"dona_tutti_api/middleware"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(g *echo.Group, service Service, rbacService middleware.RBACService) {
	handler := NewHandler(service)
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	// Payment methods routes
	paymentMethodGroup := g.Group("/payment-methods")

	// Public routes
	paymentMethodGroup.GET("", handler.GetPaymentMethods)
	paymentMethodGroup.GET("/:id", handler.GetPaymentMethod)

	// Campaign payment methods routes
	campaignPaymentGroup := g.Group("/campaigns/:campaign_id/payment-methods")
	campaignPaymentGroup.GET("", handler.GetCampaignPaymentMethods)

	// Protected routes with authentication
	authGroup := campaignPaymentGroup.Group("", middleware.RequireAuth())

	// Admin or organizer routes
	adminOrOrganizerGroup := authGroup.Group("", rbacMiddleware.Combine(
		rbacMiddleware.RequireRole("admin"),
		rbacMiddleware.RequireRole("organizer"),
	))
	adminOrOrganizerGroup.POST("", handler.CreateCampaignPaymentMethod)
	adminOrOrganizerGroup.PUT("/:id", handler.UpdateCampaignPaymentMethod)
	adminOrOrganizerGroup.DELETE("/:id", handler.DeleteCampaignPaymentMethod)
}

// @Summary Get all payment methods
// @Description Get a list of all available payment methods
// @Tags payment-methods
// @Accept json
// @Produce json
// @Success 200 {array} PaymentMethod
// @Router /payment-methods [get]
func (h *Handler) GetPaymentMethods(c echo.Context) error {
	ctx := c.Request().Context()

	paymentMethods, err := h.service.GetPaymentMethods(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get payment methods")
	}

	return c.JSON(http.StatusOK, paymentMethods)
}

// @Summary Get payment method by ID
// @Description Get a specific payment method by its ID
// @Tags payment-methods
// @Accept json
// @Produce json
// @Param id path int true "Payment Method ID"
// @Success 200 {object} PaymentMethod
// @Failure 404 {object} map[string]string
// @Router /payment-methods/{id} [get]
func (h *Handler) GetPaymentMethod(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment method ID")
	}

	paymentMethod, err := h.service.GetPaymentMethod(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Payment method not found")
	}

	return c.JSON(http.StatusOK, paymentMethod)
}

// @Summary Get campaign payment methods
// @Description Get all payment methods configured for a specific campaign
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign_id path string true "Campaign ID"
// @Success 200 {array} CampaignPaymentMethod
// @Failure 400 {object} map[string]string
// @Router /campaigns/{campaign_id}/payment-methods [get]
func (h *Handler) GetCampaignPaymentMethods(c echo.Context) error {
	ctx := c.Request().Context()

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	campaignPaymentMethods, err := h.service.GetCampaignPaymentMethods(ctx, campaignID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get campaign payment methods")
	}

	return c.JSON(http.StatusOK, campaignPaymentMethods)
}

// @Summary Create campaign payment method
// @Description Add a payment method to a campaign with optional details
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign_id path string true "Campaign ID"
// @Param request body CreateCampaignPaymentMethodRequest true "Campaign Payment Method Request"
// @Success 201 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /campaigns/{campaign_id}/payment-methods [post]
func (h *Handler) CreateCampaignPaymentMethod(c echo.Context) error {
	ctx := c.Request().Context()

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	var req CreateCampaignPaymentMethodRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Ensure the campaign ID in the URL matches the one in the request
	req.CampaignID = campaignID

	if err := c.Validate(&req); err != nil {
		c.Logger().Debugf("Validation failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := h.service.CreateCampaignPaymentMethod(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// @Summary Update campaign payment method
// @Description Update a campaign payment method and its details
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign_id path string true "Campaign ID"
// @Param id path int true "Campaign Payment Method ID"
// @Param request body CreateCampaignPaymentMethodRequest true "Campaign Payment Method Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /campaigns/{campaign_id}/payment-methods/{id} [put]
func (h *Handler) UpdateCampaignPaymentMethod(c echo.Context) error {
	ctx := c.Request().Context()

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign payment method ID")
	}

	var req CreateCampaignPaymentMethodRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Ensure the campaign ID in the URL matches the one in the request
	req.CampaignID = campaignID

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.service.UpdateCampaignPaymentMethod(ctx, id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Campaign payment method updated successfully"})
}

// @Summary Delete campaign payment method
// @Description Remove a payment method from a campaign
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign_id path string true "Campaign ID"
// @Param id path int true "Campaign Payment Method ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /campaigns/{campaign_id}/payment-methods/{id} [delete]
func (h *Handler) DeleteCampaignPaymentMethod(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign payment method ID")
	}

	err = h.service.DeleteCampaignPaymentMethod(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Campaign payment method not found")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Campaign payment method deleted successfully"})
}
