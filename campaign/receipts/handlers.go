package receipts

import (
	"dona_tutti_api/s3client"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  Service
	s3Client *s3client.Client
}

func NewHandler(service Service, s3Client *s3client.Client) *Handler {
	return &Handler{
		service:  service,
		s3Client: s3Client,
	}
}

// GetReceiptsByCampaign returns all receipts for a campaign
// @Summary Get all receipts for a campaign
// @Description Get a list of all receipts for a specific campaign
// @Tags receipts
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Success 200 {array} Receipt
// @Failure 400 {object} map[string]string
// @Router /campaigns/{campaignId}/receipts [get]
func (h *Handler) GetReceiptsByCampaign(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	receipts, err := h.service.GetReceiptsByCampaign(c.Request().Context(), campaignID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, receipts)
}

// GetReceipt returns a specific receipt
// @Summary Get a receipt by ID
// @Description Get receipt details by ID
// @Tags receipts
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Receipt ID"
// @Success 200 {object} Receipt
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /campaigns/{campaignId}/receipts/{id} [get]
func (h *Handler) GetReceipt(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid receipt ID")
	}

	receipt, err := h.service.GetReceipt(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "receipt not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, receipt)
}

// CreateReceipt creates a new receipt
// @Summary Create a new receipt
// @Description Create a new receipt for a campaign
// @Tags receipts
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param receipt body ReceiptCreateRequest true "Receipt data"
// @Success 201 {object} Receipt
// @Failure 400 {object} map[string]string
// @Router /campaigns/{campaignId}/receipts [post]
func (h *Handler) CreateReceipt(c echo.Context) error {
	campaignID, err := uuid.Parse(c.Param("campaignId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	var req ReceiptCreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	receipt, err := h.service.CreateReceipt(c.Request().Context(), campaignID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, receipt)
}

// UpdateReceipt updates an existing receipt
// @Summary Update a receipt
// @Description Update an existing receipt
// @Tags receipts
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Receipt ID"
// @Param receipt body ReceiptUpdateRequest true "Receipt update data"
// @Success 200 {object} Receipt
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /campaigns/{campaignId}/receipts/{id} [put]
func (h *Handler) UpdateReceipt(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid receipt ID")
	}

	var req ReceiptUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	receipt, err := h.service.UpdateReceipt(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "receipt not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, receipt)
}

// DeleteReceipt deletes a receipt
// @Summary Delete a receipt
// @Description Delete a receipt by ID
// @Tags receipts
// @Accept json
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Receipt ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /campaigns/{campaignId}/receipts/{id} [delete]
func (h *Handler) DeleteReceipt(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid receipt ID")
	}

	if err := h.service.DeleteReceipt(c.Request().Context(), id); err != nil {
		if err.Error() == "receipt not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// UploadReceiptDocument uploads a PDF document for a receipt
// @Summary Upload receipt document
// @Description Upload a PDF document for a receipt
// @Tags receipts
// @Accept multipart/form-data
// @Produce json
// @Param campaignId path string true "Campaign ID"
// @Param id path string true "Receipt ID"
// @Param file formData file true "PDF file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /campaigns/{campaignId}/receipts/{id}/upload [post]
func (h *Handler) UploadReceiptDocument(c echo.Context) error {
	if h.s3Client == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "File upload service is not available")
	}

	receiptID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid receipt ID")
	}

	// Check if receipt exists
	_, err = h.service.GetReceipt(c.Request().Context(), receiptID)
	if err != nil {
		if err.Error() == "receipt not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Get the file from form
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		return echo.NewHTTPError(http.StatusBadRequest, "Only PDF files are allowed")
	}

	// Validate MIME type
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file")
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "application/pdf" {
		return echo.NewHTTPError(http.StatusBadRequest, "File must be a valid PDF document")
	}

	// Reset file position
	src.Seek(0, 0)
	// Create upload request
	uploadReq := s3client.UploadRequest{
		File:         src,
		Header:       file,
		ResourceType: "receipt",
		ResourceID:   receiptID.String(),
	}

	// Upload to S3
	response, err := h.s3Client.Upload(c.Request().Context(), uploadReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update receipt with document URL
	if err := h.service.UpdateDocumentURL(c.Request().Context(), receiptID, response.URL); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update receipt with document URL")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":      "Document uploaded successfully",
		"document_url": response.URL,
	})
}
