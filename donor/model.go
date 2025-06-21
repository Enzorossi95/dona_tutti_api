package donor

import (
	"time"

	"github.com/google/uuid"
)

type DonorModel struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FirstName  string    `gorm:"column:first_name;not null"`
	LastName   string    `gorm:"column:last_name;not null"`
	IsVerified bool      `gorm:"column:is_verified;default:false"`
	Phone      string    `gorm:"column:phone"`
	Email      string    `gorm:"column:email;unique;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DonorModel) TableName() string {
	return "donors"
}

func (m DonorModel) ToEntity() Donor {
	return Donor{
		ID:         m.ID,
		FirstName:  m.FirstName,
		LastName:   m.LastName,
		IsVerified: m.IsVerified,
		Phone:      m.Phone,
		Email:      m.Email,
	}
}

func (m *DonorModel) FromEntity(entity Donor) {
	m.ID = entity.ID
	m.FirstName = entity.FirstName
	m.LastName = entity.LastName
	m.IsVerified = entity.IsVerified
	m.Phone = entity.Phone
	m.Email = entity.Email
}
