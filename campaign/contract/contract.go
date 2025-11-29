package contract

import (
	"time"

	"github.com/google/uuid"
)

// CampaignContract represents the domain entity for campaign legal contracts
type CampaignContract struct {
	ID             uuid.UUID          `json:"id"`
	CampaignID     uuid.UUID          `json:"campaign_id"`
	OrganizerID    uuid.UUID          `json:"organizer_id"`
	ContractPdfURL string             `json:"contract_pdf_url"`
	ContractHash   string             `json:"contract_hash"`
	AcceptedAt     time.Time          `json:"accepted_at"`
	Acceptance     AcceptanceMetadata `json:"acceptance_metadata"`
	CreatedAt      time.Time          `json:"created_at"`
}

// AcceptanceMetadata represents the metadata collected during contract acceptance
type AcceptanceMetadata struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

// ContractData represents the data needed to generate a contract PDF
type ContractData struct {
	CampaignID       uuid.UUID
	CampaignTitle    string
	CampaignGoal     float64
	OrganizerID      uuid.UUID
	OrganizerName    string
	OrganizerEmail   string
	OrganizerPhone   string
	OrganizerAddress string
	GeneratedAt      time.Time
}

// AcceptContractRequest represents the request to accept a contract
type AcceptContractRequest struct {
	CampaignID  uuid.UUID
	OrganizerID uuid.UUID
	IP          string
	UserAgent   string
}

// ContractProof represents the proof of contract for admin view
type ContractProof struct {
	Contract      CampaignContract `json:"contract"`
	CampaignTitle string           `json:"campaign_title"`
	OrganizerName string           `json:"organizer_name"`
}

