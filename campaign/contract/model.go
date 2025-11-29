package contract

import (
	"time"

	"github.com/google/uuid"
)

// CampaignContractModel represents the database table structure with GORM tags
type CampaignContractModel struct {
	ID                  uuid.UUID `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CampaignID          uuid.UUID `gorm:"column:campaign_id;type:uuid;not null;uniqueIndex"`
	OrganizerID         uuid.UUID `gorm:"column:organizer_id;type:uuid;not null;index"`
	ContractPdfURL      string    `gorm:"column:contract_pdf_url;type:text;not null"`
	ContractHash        string    `gorm:"column:contract_hash;type:varchar(64);not null"`
	AcceptedAt          time.Time `gorm:"column:accepted_at;not null"`
	AcceptanceIP        string    `gorm:"column:acceptance_ip;type:varchar(45);not null"`
	AcceptanceUserAgent string    `gorm:"column:acceptance_user_agent;type:text"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name for GORM
func (CampaignContractModel) TableName() string {
	return "campaign_contracts"
}

// ToEntity converts a database model to a domain entity
func (m CampaignContractModel) ToEntity() CampaignContract {
	return CampaignContract{
		ID:          m.ID,
		CampaignID:  m.CampaignID,
		OrganizerID: m.OrganizerID,
		ContractPdfURL: m.ContractPdfURL,
		ContractHash:   m.ContractHash,
		AcceptedAt:     m.AcceptedAt,
		Acceptance: AcceptanceMetadata{
			IP:        m.AcceptanceIP,
			UserAgent: m.AcceptanceUserAgent,
		},
		CreatedAt: m.CreatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *CampaignContractModel) FromEntity(entity CampaignContract) {
	m.ID = entity.ID
	m.CampaignID = entity.CampaignID
	m.OrganizerID = entity.OrganizerID
	m.ContractPdfURL = entity.ContractPdfURL
	m.ContractHash = entity.ContractHash
	m.AcceptedAt = entity.AcceptedAt
	m.AcceptanceIP = entity.Acceptance.IP
	m.AcceptanceUserAgent = entity.Acceptance.UserAgent
	m.CreatedAt = entity.CreatedAt
}

