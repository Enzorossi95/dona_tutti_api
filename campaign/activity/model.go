package activity

import (
	"time"

	"github.com/google/uuid"
)

// ActivityModel represents the database table structure with GORM tags
type ActivityModel struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CampaignID  uuid.UUID `gorm:"column:campaign_id;type:uuid;not null;index"`
	Title       string    `gorm:"column:title;not null"`
	Description string    `gorm:"column:description"`
	Date        time.Time `gorm:"column:date;not null;index"`
	Type        string    `gorm:"column:type;not null"`
	Author      string    `gorm:"column:author;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (ActivityModel) TableName() string {
	return "activities"
}

// ToEntity converts a database model to a domain entity
func (m ActivityModel) ToEntity() Activity {
	return Activity{
		ID:          m.ID,
		CampaignID:  m.CampaignID,
		Title:       m.Title,
		Description: m.Description,
		Date:        m.Date,
		Type:        m.Type,
		Author:      m.Author,
		CreatedAt:   m.CreatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *ActivityModel) FromEntity(entity Activity) {
	m.ID = entity.ID
	m.CampaignID = entity.CampaignID
	m.Title = entity.Title
	m.Description = entity.Description
	m.Date = entity.Date
	m.Type = entity.Type
	m.Author = entity.Author
	m.CreatedAt = entity.CreatedAt
}