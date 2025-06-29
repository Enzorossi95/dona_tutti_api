package donation

import (
	"time"

	"github.com/google/uuid"
)

type DonationStatus string

const (
	DonationStatusCompleted DonationStatus = "completed"
	DonationStatusPending   DonationStatus = "pending"
	DonationStatusFailed    DonationStatus = "failed"
	DonationStatusRefunded  DonationStatus = "refunded"
)

// PaymentMethodInfo represents payment method information in donation context
type PaymentMethodInfo struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Donation struct {
	ID            uuid.UUID          `json:"id"`
	CampaignID    uuid.UUID          `json:"campaign_id"`
	Amount        float64            `json:"amount"`
	DonorID       uuid.UUID          `json:"donor_id"`
	Date          time.Time          `json:"date"`
	Message       *string            `json:"message,omitempty"`
	IsAnonymous   bool               `json:"is_anonymous"`
	PaymentMethodID int              `json:"payment_method_id"`
	PaymentMethod *PaymentMethodInfo `json:"payment_method,omitempty"`
	Status        DonationStatus     `json:"status"`
}
