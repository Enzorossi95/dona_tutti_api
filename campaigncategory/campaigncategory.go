package campaigncategory

import (
	"time"

	"github.com/google/uuid"
)

// CampaignCategory represents the domain entity for campaign categories
type CampaignCategory struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
