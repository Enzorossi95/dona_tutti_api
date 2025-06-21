package campaign

import (
	"time"

	"github.com/google/uuid"
)

// Campaign represents the domain entity for campaigns
type Campaign struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Goal        float64   `json:"goal"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	CategoryId  uuid.UUID `json:"category"`
	Urgency     int       `json:"urgency"`
	OrganizerId uuid.UUID `json:"organizer"`
	Status      string    `json:"status"`
}

type Summary struct {
	TotalCampaigns    int64   `json:"total_campaigns"`
	TotalGoal         float64 `json:"total_goal"`
	TotalContributors int64   `json:"total_contributors"`
}
