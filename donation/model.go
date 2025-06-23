package donation

import (
	"dona_tutti_api/donor"
	"time"

	"github.com/google/uuid"
)

type DonationModel struct {
	ID            uuid.UUID        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CampaignID    uuid.UUID        `gorm:"column:campaign_id;type:uuid;not null"`
	Amount        float64          `gorm:"column:amount;not null"`
	DonorID       uuid.UUID        `gorm:"column:donor_id;type:uuid;not null"`
	Date          time.Time        `gorm:"column:date;not null"`
	Message       *string          `gorm:"column:message"`
	IsAnonymous   bool             `gorm:"column:is_anonymous"`
	PaymentMethod PaymentMethod    `gorm:"column:payment_method;type:varchar(20);not null"`
	Status        DonationStatus   `gorm:"column:status;type:varchar(20);not null"`
	CreatedAt     time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time        `gorm:"column:updated_at;autoUpdateTime"`
	Donor         donor.DonorModel `gorm:"foreignKey:DonorID"`
}

func (DonationModel) TableName() string {
	return "donations"
}

func (m DonationModel) ToEntity() Donation {
	return Donation{
		ID:            m.ID,
		CampaignID:    m.CampaignID,
		Amount:        m.Amount,
		DonorID:       m.DonorID,
		Date:          m.Date,
		Message:       m.Message,
		IsAnonymous:   m.IsAnonymous,
		PaymentMethod: m.PaymentMethod,
		Status:        m.Status,
	}
}

func (m *DonationModel) FromEntity(entity Donation) {
	m.ID = entity.ID
	m.CampaignID = entity.CampaignID
	m.Amount = entity.Amount
	m.DonorID = entity.DonorID
	m.Date = entity.Date
	m.Message = entity.Message
	m.IsAnonymous = entity.IsAnonymous
	m.PaymentMethod = entity.PaymentMethod
	m.Status = entity.Status
}
