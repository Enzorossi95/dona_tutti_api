package receipts

import (
	"time"

	"github.com/google/uuid"
)

// Receipt represents the domain entity for campaign receipts
type Receipt struct {
	ID          uuid.UUID  `json:"id"`
	CampaignID  uuid.UUID  `json:"campaign_id"`
	Provider    string     `json:"provider"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Total       float64    `json:"total"`
	Quantity    int        `json:"quantity"`
	Date        time.Time  `json:"date"`
	DocumentURL *string    `json:"document_url,omitempty"`
	Note        *string    `json:"note,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ReceiptCreateRequest represents the request to create a new receipt
type ReceiptCreateRequest struct {
	Provider    string    `json:"provider" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Total       float64   `json:"total" validate:"required,gt=0"`
	Quantity    int       `json:"quantity" validate:"omitempty,gte=1"`
	Date        time.Time `json:"date" validate:"required"`
	Note        *string   `json:"note,omitempty"`
}

// ReceiptUpdateRequest represents a partial update request for receipts
type ReceiptUpdateRequest struct {
	Provider    *string    `json:"provider,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Total       *float64   `json:"total,omitempty" validate:"omitempty,gt=0"`
	Quantity    *int       `json:"quantity,omitempty" validate:"omitempty,gte=1"`
	Date        *time.Time `json:"date,omitempty"`
	Note        *string    `json:"note,omitempty"`
}