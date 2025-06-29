package campaign

import (
	"time"

	"github.com/google/uuid"
)

type CampaignRequest struct {
	ID                uuid.UUID `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Image             string    `json:"image"`
	Goal              float64   `json:"goal"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	Location          string    `json:"location"`
	CategoryId        uuid.UUID `json:"category"`
	Urgency           int       `json:"urgency"`
	OrganizerId       uuid.UUID `json:"organizer"`
	Status            string    `json:"status"`
	PaymentMethodsIds []int     `json:"payment_methods_ids,omitempty"`
}
