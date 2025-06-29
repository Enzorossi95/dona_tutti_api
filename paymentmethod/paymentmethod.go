package paymentmethod

import (
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents the domain entity for payment methods
type PaymentMethod struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// CampaignPaymentMethod represents the domain entity for campaign payment method associations
type CampaignPaymentMethod struct {
	ID              int              `json:"id"`
	CampaignID      uuid.UUID        `json:"campaign_id"`
	PaymentMethodID int              `json:"payment_method_id"`
	Instructions    *string          `json:"instructions,omitempty"`
	IsActive        bool             `json:"is_active"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	PaymentMethod   *PaymentMethod   `json:"payment_method,omitempty"`
	TransferDetails []TransferDetail `json:"transfer_details,omitempty"`
	CashLocations   []CashLocation   `json:"cash_locations,omitempty"`
}

// TransferDetail represents the domain entity for bank transfer details
type TransferDetail struct {
	ID                      int       `json:"id"`
	CampaignPaymentMethodID int       `json:"campaign_payment_method_id"`
	BankName                string    `json:"bank_name"`
	AccountHolder           string    `json:"account_holder"`
	CBU                     string    `json:"cbu"`
	Alias                   *string   `json:"alias,omitempty"`
	SwiftCode               *string   `json:"swift_code,omitempty"`
	AdditionalNotes         *string   `json:"additional_notes,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// CashLocation represents the domain entity for cash payment locations
type CashLocation struct {
	ID                      int       `json:"id"`
	CampaignPaymentMethodID int       `json:"campaign_payment_method_id"`
	LocationName            string    `json:"location_name"`
	Address                 string    `json:"address"`
	ContactInfo             *string   `json:"contact_info,omitempty"`
	AvailableHours          *string   `json:"available_hours,omitempty"`
	AdditionalNotes         *string   `json:"additional_notes,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// CreateCampaignPaymentMethodRequest represents the request to create a campaign payment method
type CreateCampaignPaymentMethodRequest struct {
	CampaignID      uuid.UUID             `json:"campaign_id" validate:"required"`
	PaymentMethodID int                   `json:"payment_method_id" validate:"required"`
	Instructions    *string               `json:"instructions,omitempty"`
	TransferDetails *CreateTransferDetail `json:"transfer_details,omitempty"`
	CashLocations   []CreateCashLocation  `json:"cash_locations,omitempty"`
}

// CreateTransferDetail represents the request to create transfer details
type CreateTransferDetail struct {
	BankName        string  `json:"bank_name" validate:"required"`
	AccountHolder   string  `json:"account_holder" validate:"required"`
	CBU             string  `json:"cbu" validate:"required"`
	Alias           *string `json:"alias,omitempty"`
	SwiftCode       *string `json:"swift_code,omitempty"`
	AdditionalNotes *string `json:"additional_notes,omitempty"`
}

// CreateCashLocation represents the request to create cash location
type CreateCashLocation struct {
	LocationName    string  `json:"location_name" validate:"required"`
	Address         string  `json:"address" validate:"required"`
	ContactInfo     *string `json:"contact_info,omitempty"`
	AvailableHours  *string `json:"available_hours,omitempty"`
	AdditionalNotes *string `json:"additional_notes,omitempty"`
}
