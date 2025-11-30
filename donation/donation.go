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

// IsValidStatus checks if a donation status is valid
func IsValidStatus(status DonationStatus) bool {
	switch status {
	case DonationStatusCompleted, DonationStatusPending, DonationStatusFailed, DonationStatusRefunded:
		return true
	default:
		return false
	}
}

// PaymentMethodInfo represents payment method information in donation context
type PaymentMethodInfo struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// DonorResponse represents donor information in donation responses
type DonorResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
}

type DonorInfo struct {
	Name     string  `json:"name" validate:"required"`
	LastName string  `json:"last_name" validate:"required"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type CreateDonationRequest struct {
	Amount          float64    `json:"amount"`
	Message         *string    `json:"message,omitempty"`
	IsAnonymous     bool       `json:"is_anonymous"`
	PaymentMethodID int        `json:"payment_method_id"`
	DonorID         *uuid.UUID `json:"donor_id,omitempty"`
	Donor           *DonorInfo `json:"donor,omitempty"`
}

type UpdateDonationStatusRequest struct {
	Status DonationStatus `json:"status" validate:"required"`
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
	Donor         *DonorResponse     `json:"donor,omitempty"`
	Status        DonationStatus     `json:"status"`
	ReceiptURL    *string            `json:"receipt_url,omitempty"`
}
