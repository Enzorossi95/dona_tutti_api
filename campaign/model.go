package campaign

import (
	"time"

	"github.com/google/uuid"
)

// CampaignModel represents the database table structure with GORM tags
type CampaignModel struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Title       string    `gorm:"column:title;not null"`
	Description string    `gorm:"column:description;not null"`
	Image       string    `gorm:"column:image"`
	Goal        float64   `gorm:"column:goal;not null;check:goal > 0"`
	StartDate   time.Time `gorm:"column:start_date;not null"`
	EndDate     time.Time `gorm:"column:end_date;not null"`
	Location    string    `gorm:"column:location"`
	Urgency     int       `gorm:"column:urgency;check:urgency >= 1 AND urgency <= 10"`
	Status      string    `gorm:"column:status;default:active"`
	CategoryID  uuid.UUID `gorm:"column:category_id"`
	OrganizerID uuid.UUID `gorm:"column:organizer_id"`
}

// TableName specifies the table name for GORM
func (CampaignModel) TableName() string {
	return "campaigns"
}

// ToEntity converts a database model to a domain entity
func (m CampaignModel) ToEntity() Campaign {
	return Campaign{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		Title:       m.Title,
		Description: m.Description,
		Image:       m.Image,
		Goal:        m.Goal,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		Location:    m.Location,
		CategoryId:  m.CategoryID,
		Urgency:     m.Urgency,
		OrganizerId: m.OrganizerID,
		Status:      m.Status,
	}
}

// FromEntity converts a domain entity to a database model
func (m *CampaignModel) FromEntity(entity Campaign) {
	m.ID = entity.ID
	m.CreatedAt = entity.CreatedAt
	m.Title = entity.Title
	m.Description = entity.Description
	m.Image = entity.Image
	m.Goal = entity.Goal
	m.StartDate = entity.StartDate
	m.EndDate = entity.EndDate
	m.Location = entity.Location
	m.Urgency = entity.Urgency
	m.Status = entity.Status
	m.CategoryID = entity.CategoryId
	m.OrganizerID = entity.OrganizerId
}
