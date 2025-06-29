package paymentmethod

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethodModel struct {
	ID        int       `gorm:"primaryKey;column:id;autoIncrement"`
	Code      string    `gorm:"column:code;uniqueIndex;not null;size:30"`
	Name      string    `gorm:"column:name;not null;size:50"`
	IsActive  bool      `gorm:"column:is_active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (PaymentMethodModel) TableName() string {
	return "payment_methods"
}

// ToEntity converts a database model to a domain entity
func (m PaymentMethodModel) ToEntity() PaymentMethod {
	return PaymentMethod{
		ID:       m.ID,
		Code:     m.Code,
		Name:     m.Name,
		IsActive: m.IsActive,
	}
}

// FromEntity converts a domain entity to a database model
func (m *PaymentMethodModel) FromEntity(entity PaymentMethod) {
	m.ID = entity.ID
	m.Code = entity.Code
	m.Name = entity.Name
	m.IsActive = entity.IsActive
}

type CampaignPaymentMethodModel struct {
	ID              int                   `gorm:"primaryKey;column:id;autoIncrement"`
	CampaignID      uuid.UUID             `gorm:"column:campaign_id;type:uuid;not null"`
	PaymentMethodID int                   `gorm:"column:payment_method_id;not null"`
	Instructions    *string               `gorm:"column:instructions;type:text"`
	IsActive        bool                  `gorm:"column:is_active;default:true"`
	CreatedAt       time.Time             `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time             `gorm:"column:updated_at;autoUpdateTime"`
	PaymentMethod   PaymentMethodModel    `gorm:"foreignKey:PaymentMethodID;constraint:OnDelete:CASCADE"`
	TransferDetails []TransferDetailModel `gorm:"foreignKey:CampaignPaymentMethodID"`
	CashLocations   []CashLocationModel   `gorm:"foreignKey:CampaignPaymentMethodID"`
}

// TableName specifies the table name for GORM
func (CampaignPaymentMethodModel) TableName() string {
	return "campaign_payment_methods"
}

// ToEntity converts a database model to a domain entity
func (m CampaignPaymentMethodModel) ToEntity() CampaignPaymentMethod {
	entity := CampaignPaymentMethod{
		ID:              m.ID,
		CampaignID:      m.CampaignID,
		PaymentMethodID: m.PaymentMethodID,
		Instructions:    m.Instructions,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	if m.PaymentMethod.ID != 0 {
		paymentMethod := m.PaymentMethod.ToEntity()
		entity.PaymentMethod = &paymentMethod
	}

	if len(m.TransferDetails) > 0 {
		entity.TransferDetails = make([]TransferDetail, len(m.TransferDetails))
		for i, transferDetail := range m.TransferDetails {
			entity.TransferDetails[i] = transferDetail.ToEntity()
		}
	}

	if len(m.CashLocations) > 0 {
		entity.CashLocations = make([]CashLocation, len(m.CashLocations))
		for i, location := range m.CashLocations {
			entity.CashLocations[i] = location.ToEntity()
		}
	}

	return entity
}

// FromEntity converts a domain entity to a database model
func (m *CampaignPaymentMethodModel) FromEntity(entity CampaignPaymentMethod) {
	m.ID = entity.ID
	m.CampaignID = entity.CampaignID
	m.PaymentMethodID = entity.PaymentMethodID
	m.Instructions = entity.Instructions
	m.IsActive = entity.IsActive
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}

// TransferDetailModel represents the database table structure with GORM tags
type TransferDetailModel struct {
	ID                      int       `gorm:"primaryKey;column:id;autoIncrement"`
	CampaignPaymentMethodID int       `gorm:"column:campaign_payment_method_id;not null"`
	BankName                string    `gorm:"column:bank_name;not null;size:100"`
	AccountHolder           string    `gorm:"column:account_holder;not null;size:100"`
	CBU                     string    `gorm:"column:cbu;not null;size:22"`
	Alias                   *string   `gorm:"column:alias;size:30"`
	SwiftCode               *string   `gorm:"column:swift_code;size:11"`
	AdditionalNotes         *string   `gorm:"column:additional_notes;type:text"`
	CreatedAt               time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt               time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (TransferDetailModel) TableName() string {
	return "transfer_details"
}

// ToEntity converts a database model to a domain entity
func (m TransferDetailModel) ToEntity() TransferDetail {
	return TransferDetail{
		ID:                      m.ID,
		CampaignPaymentMethodID: m.CampaignPaymentMethodID,
		BankName:                m.BankName,
		AccountHolder:           m.AccountHolder,
		CBU:                     m.CBU,
		Alias:                   m.Alias,
		SwiftCode:               m.SwiftCode,
		AdditionalNotes:         m.AdditionalNotes,
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *TransferDetailModel) FromEntity(entity TransferDetail) {
	m.ID = entity.ID
	m.CampaignPaymentMethodID = entity.CampaignPaymentMethodID
	m.BankName = entity.BankName
	m.AccountHolder = entity.AccountHolder
	m.CBU = entity.CBU
	m.Alias = entity.Alias
	m.SwiftCode = entity.SwiftCode
	m.AdditionalNotes = entity.AdditionalNotes
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}

// CashLocationModel represents the database table structure with GORM tags
type CashLocationModel struct {
	ID                      int       `gorm:"primaryKey;column:id;autoIncrement"`
	CampaignPaymentMethodID int       `gorm:"column:campaign_payment_method_id;not null"`
	LocationName            string    `gorm:"column:location_name;not null;size:100"`
	Address                 string    `gorm:"column:address;not null;size:200"`
	ContactInfo             *string   `gorm:"column:contact_info;size:100"`
	AvailableHours          *string   `gorm:"column:available_hours;size:100"`
	AdditionalNotes         *string   `gorm:"column:additional_notes;type:text"`
	CreatedAt               time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt               time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (CashLocationModel) TableName() string {
	return "cash_locations"
}

// ToEntity converts a database model to a domain entity
func (m CashLocationModel) ToEntity() CashLocation {
	return CashLocation{
		ID:                      m.ID,
		CampaignPaymentMethodID: m.CampaignPaymentMethodID,
		LocationName:            m.LocationName,
		Address:                 m.Address,
		ContactInfo:             m.ContactInfo,
		AvailableHours:          m.AvailableHours,
		AdditionalNotes:         m.AdditionalNotes,
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *CashLocationModel) FromEntity(entity CashLocation) {
	m.ID = entity.ID
	m.CampaignPaymentMethodID = entity.CampaignPaymentMethodID
	m.LocationName = entity.LocationName
	m.Address = entity.Address
	m.ContactInfo = entity.ContactInfo
	m.AvailableHours = entity.AvailableHours
	m.AdditionalNotes = entity.AdditionalNotes
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}
