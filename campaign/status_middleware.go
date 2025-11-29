package campaign

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ContractChecker defines the interface for checking contract existence
type ContractChecker interface {
	HasContract(campaignID uuid.UUID) (bool, error)
}

// StatusValidationMiddleware provides middleware for validating campaign status transitions
type StatusValidationMiddleware struct {
	contractChecker ContractChecker
}

// NewStatusValidationMiddleware creates a new status validation middleware
func NewStatusValidationMiddleware(contractChecker ContractChecker) *StatusValidationMiddleware {
	return &StatusValidationMiddleware{
		contractChecker: contractChecker,
	}
}

// ValidateStatusTransition is a middleware that validates status transitions
func (m *StatusValidationMiddleware) ValidateStatusTransition() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// This middleware can be used to add additional validation logic
			// before allowing status transitions
			// For now, the validation is handled in the service layer
			return next(c)
		}
	}
}

// RequireContractForApproval is a middleware that ensures a contract exists before approval
func (m *StatusValidationMiddleware) RequireContractForApproval() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Parse request to check if status is being changed to active
			var req map[string]interface{}
			if err := c.Bind(&req); err != nil {
				return next(c) // Let the handler deal with invalid request
			}

			// Check if status is being updated to active
			status, ok := req["status"].(string)
			if !ok || status != StatusActive {
				return next(c) // Not updating to active, continue
			}

			// Get campaign ID from URL
			campaignIDStr := c.Param("id")
			campaignID, err := uuid.Parse(campaignIDStr)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error": "Invalid campaign ID format",
				})
			}

			// Check if contract exists
			hasContract, err := m.contractChecker.HasContract(campaignID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "Failed to check contract status",
				})
			}

			if !hasContract {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error": "Campaign must have a signed contract before activation",
				})
			}

			return next(c)
		}
	}
}

// ValidateStatusValue is a middleware that validates if the status value is valid
func ValidateStatusValue() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Parse request to check status value
			var req map[string]interface{}
			if err := c.Bind(&req); err != nil {
				return next(c) // Let the handler deal with invalid request
			}

			// Check if status field exists
			status, ok := req["status"].(string)
			if !ok {
				return next(c) // No status field, continue
			}

			// Validate status value
			if !IsValidStatus(status) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error": "Invalid status value. Valid statuses: draft, pending_approval, active, paused, completed, rejected",
				})
			}

			return next(c)
		}
	}
}

