package campaign

import (
	"time"

	"dona_tutti_api/organizer"

	"github.com/google/uuid"
)

// CampaignPaymentMethod represents a payment method available for a campaign
type CampaignPaymentMethod struct {
	ID              int     `json:"id"`
	PaymentMethodID int     `json:"payment_method_id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Instructions    *string `json:"instructions,omitempty"`
}

// Campaign represents the domain entity for campaigns
type Campaign struct {
	ID               uuid.UUID               `json:"id"`
	CreatedAt        time.Time               `json:"created_at"`
	Title            string                  `json:"title"`
	Description      string                  `json:"description"`
	Image            string                  `json:"image"`
	Goal             float64                 `json:"goal"`
	StartDate        time.Time               `json:"start_date"`
	EndDate          time.Time               `json:"end_date"`
	Location         string                  `json:"location"`
	CategoryId       uuid.UUID               `json:"category"`
	Urgency          int                     `json:"urgency"`
	Organizer        *organizer.Organizer    `json:"organizer"`
	Status           string                  `json:"status"`
	PaymentMethods   []CampaignPaymentMethod `json:"payment_methods,omitempty"`
	BeneficiaryName  *string                 `json:"beneficiary_name,omitempty"`
	BeneficiaryAge   *int                    `json:"beneficiary_age,omitempty"`
	CurrentSituation *string                 `json:"current_situation,omitempty"`
	UrgencyReason    *string                 `json:"urgency_reason,omitempty"`
}

type Summary struct {
	TotalCampaigns    int64   `json:"total_campaigns"`
	TotalGoal         float64 `json:"total_goal"`
	TotalContributors int64   `json:"total_contributors"`
}
