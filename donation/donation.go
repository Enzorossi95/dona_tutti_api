package donation

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string
type DonationStatus string

const (
	PaymentMethodMercadoPago   PaymentMethod = "MercadoPago"
	PaymentMethodTransferencia PaymentMethod = "Transferencia"
	PaymentMethodEfectivo      PaymentMethod = "Efectivo"
)

const (
	DonationStatusCompleted DonationStatus = "completed"
	DonationStatusPending   DonationStatus = "pending"
	DonationStatusFailed    DonationStatus = "failed"
	DonationStatusRefunded  DonationStatus = "refunded"
)

type Donation struct {
	ID            uuid.UUID      `json:"id"`
	CampaignID    uuid.UUID      `json:"campaign_id"`
	Amount        float64        `json:"amount"`
	DonorID       uuid.UUID      `json:"donor_id"`
	Date          time.Time      `json:"date"`
	Message       *string        `json:"message,omitempty"`
	IsAnonymous   bool           `json:"is_anonymous"`
	PaymentMethod PaymentMethod  `json:"payment_method"`
	Status        DonationStatus `json:"status"`
}
