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
	OrganizerID      uuid.UUID               `json:"organizer_id"`
	Status           string                  `json:"status"`
	ContractSigned   bool                    `json:"contract_signed"`
	PaymentMethods   []CampaignPaymentMethod `json:"payment_methods,omitempty"`
	BeneficiaryName  *string                 `json:"beneficiary_name,omitempty"`
	BeneficiaryAge   *int                    `json:"beneficiary_age,omitempty"`
	CurrentSituation *string                 `json:"current_situation,omitempty"`
	UrgencyReason    *string                 `json:"urgency_reason,omitempty"`
}

// Valid campaign statuses
const (
	StatusDraft           = "draft"
	StatusPendingApproval = "pending_approval"
	StatusActive          = "active"
	StatusPaused          = "paused"
	StatusCompleted       = "completed"
	StatusRejected        = "rejected"
)

// ValidStatuses returns all valid campaign statuses
func ValidStatuses() []string {
	return []string{
		StatusDraft,
		StatusPendingApproval,
		StatusActive,
		StatusPaused,
		StatusCompleted,
		StatusRejected,
	}
}

// IsValidStatus checks if a status is valid
func IsValidStatus(status string) bool {
	validStatuses := map[string]bool{
		StatusDraft:           true,
		StatusPendingApproval: true,
		StatusActive:          true,
		StatusPaused:          true,
		StatusCompleted:       true,
		StatusRejected:        true,
	}
	return validStatuses[status]
}

// CanTransitionTo checks if a status transition is valid
func CanTransitionTo(from, to string) bool {
	validTransitions := map[string][]string{
		StatusDraft:           {StatusPendingApproval, StatusRejected},
		StatusPendingApproval: {StatusActive, StatusRejected},
		StatusActive:          {StatusPaused, StatusCompleted},
		StatusPaused:          {StatusActive, StatusCompleted},
		StatusCompleted:       {}, // Terminal state
		StatusRejected:        {}, // Terminal state
	}

	allowedTransitions, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == to {
			return true
		}
	}
	return false
}

type Summary struct {
	TotalCampaigns    int64   `json:"total_campaigns"`
	TotalGoal         float64 `json:"total_goal"`
	TotalContributors int64   `json:"total_contributors"`
}
