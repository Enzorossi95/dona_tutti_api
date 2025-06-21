package campaigncategory

import (
	"time"

	"github.com/google/uuid"
)

// CampaignCategoryModel represents the database table structure with GORM tags
type CampaignCategoryModel struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	Name        string    `gorm:"column:name;not null;unique"`
	Description string    `gorm:"column:description"`
}

// TableName specifies the table name for GORM
func (CampaignCategoryModel) TableName() string {
	return "campaign_categories"
}

// ToEntity converts a database model to a domain entity
func (m CampaignCategoryModel) ToEntity() CampaignCategory {
	return CampaignCategory{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		Name:        m.Name,
		Description: m.Description,
	}
}

// FromEntity converts a domain entity to a database model
func (m *CampaignCategoryModel) FromEntity(entity CampaignCategory) {
	m.ID = entity.ID
	m.CreatedAt = entity.CreatedAt
	m.Name = entity.Name
	m.Description = entity.Description
}
