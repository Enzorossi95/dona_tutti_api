package contract

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for campaign contracts
type Handler struct {
	service Service
}

// NewHandler creates a new contract handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the contract routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	contracts := g.Group("/campaigns/:id/contract")
	contracts.POST("/generate", h.GenerateContract)
	contracts.GET("", h.GetContract)
	contracts.POST("/accept", h.AcceptContract)
	contracts.GET("/proof", h.GetContractProof)
}

// AcceptContractRequestDTO represents the request to accept a contract
type AcceptContractRequestDTO struct {
	OrganizerID uuid.UUID `json:"organizer_id"`
}

// GenerateContract handles POST /api/campaigns/:id/contract/generate
// @Summary Generate contract PDF for a campaign
// @Description Generates a legal contract PDF using campaign data from database
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} map[string]interface{} "Contract generated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Campaign not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/campaigns/{id}/contract/generate [post]
func (h *Handler) GenerateContract(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Generate contract (service fetches all data from database)
	url, err := h.service.GenerateContract(c.Request().Context(), campaignID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Contract generated successfully",
		"contract_url": url,
	})
}

// GetContract handles GET /api/campaigns/:id/contract
// @Summary Get contract for a campaign
// @Description Retrieves the contract information for a campaign
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} CampaignContract "Contract information"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Contract not found"
// @Router /api/campaigns/{id}/contract [get]
func (h *Handler) GetContract(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Get contract
	contract, err := h.service.GetContract(c.Request().Context(), campaignID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Contract not found",
		})
	}

	return c.JSON(http.StatusOK, contract)
}

// AcceptContract handles POST /api/campaigns/:id/contract/accept
// @Summary Accept a contract
// @Description Records the acceptance of a contract with digital signature metadata
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Param request body AcceptContractRequestDTO true "Contract acceptance request"
// @Success 200 {object} map[string]interface{} "Contract accepted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/campaigns/{id}/contract/accept [post]
func (h *Handler) AcceptContract(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Parse request body
	var reqDTO AcceptContractRequestDTO
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// Get IP address from request
	ip := c.RealIP()

	// Get user agent from request
	userAgent := c.Request().UserAgent()

	// Create acceptance request
	req := AcceptContractRequest{
		CampaignID:  campaignID,
		OrganizerID: reqDTO.OrganizerID,
		IP:          ip,
		UserAgent:   userAgent,
	}

	// Accept contract
	if err := h.service.AcceptContract(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Contract accepted successfully",
		"status":  "pending_approval",
	})
}

// GetContractProof handles GET /api/campaigns/:id/contract/proof
// @Summary Get contract proof for admin
// @Description Retrieves the contract proof with full details for admin review
// @Tags contracts
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} ContractProof "Contract proof"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Contract not found"
// @Router /api/campaigns/{id}/contract/proof [get]
func (h *Handler) GetContractProof(c echo.Context) error {
	// Parse campaign ID from URL
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid campaign ID format",
		})
	}

	// Get contract proof
	proof, err := h.service.GetContractProof(c.Request().Context(), campaignID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Contract not found",
		})
	}

	return c.JSON(http.StatusOK, proof)
}
