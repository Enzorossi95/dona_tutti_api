package receipts

import (
	"time"

	"github.com/google/uuid"
)

// ReceiptModel represents the database table structure with GORM tags
type ReceiptModel struct {
	ID          uuid.UUID  `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CampaignID  uuid.UUID  `gorm:"column:campaign_id;type:uuid;not null;index"`
	Provider    string     `gorm:"column:provider;type:varchar(255);not null"`
	Name        string     `gorm:"column:name;type:varchar(255);not null"`
	Description string     `gorm:"column:description;type:text"`
	Total       float64    `gorm:"column:total;type:decimal(12,2);not null"`
	Quantity    int        `gorm:"column:quantity;type:integer;default:1"`
	Date        time.Time  `gorm:"column:date;not null;index"`
	DocumentURL *string    `gorm:"column:document_url;type:varchar(500)"`
	Note        *string    `gorm:"column:note;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (ReceiptModel) TableName() string {
	return "receipts"
}

// ToEntity converts a database model to a domain entity
func (m ReceiptModel) ToEntity() Receipt {
	return Receipt{
		ID:          m.ID,
		CampaignID:  m.CampaignID,
		Provider:    m.Provider,
		Name:        m.Name,
		Description: m.Description,
		Total:       m.Total,
		Quantity:    m.Quantity,
		Date:        m.Date,
		DocumentURL: m.DocumentURL,
		Note:        m.Note,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *ReceiptModel) FromEntity(entity Receipt) {
	m.ID = entity.ID
	m.CampaignID = entity.CampaignID
	m.Provider = entity.Provider
	m.Name = entity.Name
	m.Description = entity.Description
	m.Total = entity.Total
	m.Quantity = entity.Quantity
	m.Date = entity.Date
	m.DocumentURL = entity.DocumentURL
	m.Note = entity.Note
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}