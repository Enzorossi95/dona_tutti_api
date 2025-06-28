package organizer

import (
	"time"

	"github.com/google/uuid"
)

// OrganizerModel represents the database table structure with GORM tags
type OrganizerModel struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null;index"`
	Name      string    `gorm:"column:name;not null"`
	Avatar    string    `gorm:"column:avatar"`
	Verified  bool      `gorm:"column:verified;default:false"`
}

// TableName specifies the table name for GORM
func (OrganizerModel) TableName() string {
	return "organizers"
}

// ToEntity converts a database model to a domain entity
func (m OrganizerModel) ToEntity() Organizer {
	return Organizer{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UserID:    m.UserID,
		Name:      m.Name,
		Avatar:    m.Avatar,
		Verified:  m.Verified,
	}
}

// FromEntity converts a domain entity to a database model
func (m *OrganizerModel) FromEntity(entity Organizer) {
	m.ID = entity.ID
	m.CreatedAt = entity.CreatedAt
	m.UserID = entity.UserID
	m.Name = entity.Name
	m.Avatar = entity.Avatar
	m.Verified = entity.Verified
}
