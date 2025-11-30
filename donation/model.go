package donation

import (
	"dona_tutti_api/donor"
	"time"

	"github.com/google/uuid"
)

// PaymentMethodModel represents payment method info for donations
type PaymentMethodModel struct {
	ID   int    `gorm:"column:id"`
	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`
}

type DonationModel struct {
	ID              uuid.UUID           `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CampaignID      uuid.UUID           `gorm:"column:campaign_id;type:uuid;not null"`
	Amount          float64             `gorm:"column:amount;not null"`
	DonorID         uuid.UUID           `gorm:"column:donor_id;type:uuid;not null"`
	Date            time.Time           `gorm:"column:date;not null"`
	Message         *string             `gorm:"column:message"`
	IsAnonymous     bool                `gorm:"column:is_anonymous"`
	PaymentMethodID int                 `gorm:"column:payment_method_id;not null"`
	Status          DonationStatus      `gorm:"column:status;type:varchar(20);not null"`
	ReceiptURL      *string             `gorm:"column:receipt_url;type:varchar(500)"`
	CreatedAt       time.Time           `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time           `gorm:"column:updated_at;autoUpdateTime"`
	Donor           donor.DonorModel    `gorm:"foreignKey:DonorID"`
	PaymentMethod   *PaymentMethodModel `gorm:"-"`
}

func (DonationModel) TableName() string {
	return "donations"
}

func (m DonationModel) ToEntity() Donation {
	donation := Donation{
		ID:              m.ID,
		CampaignID:      m.CampaignID,
		Amount:          m.Amount,
		DonorID:         m.DonorID,
		Date:            m.Date,
		Message:         m.Message,
		IsAnonymous:     m.IsAnonymous,
		PaymentMethodID: m.PaymentMethodID,
		Status:          m.Status,
		ReceiptURL:      m.ReceiptURL,
	}

	// Convert payment method info if available
	if m.PaymentMethod != nil {
		donation.PaymentMethod = &PaymentMethodInfo{
			ID:   m.PaymentMethod.ID,
			Code: m.PaymentMethod.Code,
			Name: m.PaymentMethod.Name,
		}
	}

	// Convert donor info if available and not anonymous
	if !m.IsAnonymous && m.Donor.ID != uuid.Nil {
		donation.Donor = &DonorResponse{
			ID:        m.Donor.ID,
			FirstName: m.Donor.FirstName,
			LastName:  m.Donor.LastName,
			Email:     m.Donor.Email,
			Phone:     m.Donor.Phone,
		}
	}

	return donation
}

func (m *DonationModel) FromEntity(entity Donation) {
	m.ID = entity.ID
	m.CampaignID = entity.CampaignID
	m.Amount = entity.Amount
	m.DonorID = entity.DonorID
	m.Date = entity.Date
	m.Message = entity.Message
	m.IsAnonymous = entity.IsAnonymous
	m.PaymentMethodID = entity.PaymentMethodID
	m.Status = entity.Status
	m.ReceiptURL = entity.ReceiptURL
}
