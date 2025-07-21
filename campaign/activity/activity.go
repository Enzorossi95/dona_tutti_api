package activity

import (
	"time"

	"github.com/google/uuid"
)

// Activity represents the domain entity for campaign activities
type Activity struct {
	ID          uuid.UUID `json:"id"`
	CampaignID  uuid.UUID `json:"campaign_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
}

// ActivityUpdateRequest represents a partial update request for activities
type ActivityUpdateRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	Type        *string    `json:"type,omitempty"`
	Author      *string    `json:"author,omitempty"`
}