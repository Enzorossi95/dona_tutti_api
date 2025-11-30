package closure

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TransparencyBreakdownJSON for JSONB in PostgreSQL
type TransparencyBreakdownJSON TransparencyBreakdown

// Value implements the driver.Valuer interface
func (t TransparencyBreakdownJSON) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan implements the sql.Scanner interface
func (t *TransparencyBreakdownJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, t)
}

// CampaignClosureReportModel represents the database table structure
type CampaignClosureReportModel struct {
	ID                    uuid.UUID                 `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CampaignID            uuid.UUID                 `gorm:"column:campaign_id;type:uuid;not null;uniqueIndex"`
	ClosureType           string                    `gorm:"column:closure_type;type:varchar(50);not null"`
	ClosureReason         *string                   `gorm:"column:closure_reason;type:text"`
	ClosedBy              *uuid.UUID                `gorm:"column:closed_by;type:uuid"`
	TotalRaised           float64                   `gorm:"column:total_raised;type:decimal(12,2);not null;default:0"`
	TotalDonors           int                       `gorm:"column:total_donors;not null;default:0"`
	TotalDonations        int                       `gorm:"column:total_donations;not null;default:0"`
	CampaignGoal          float64                   `gorm:"column:campaign_goal;type:decimal(12,2);not null"`
	GoalPercentage        float64                   `gorm:"column:goal_percentage;type:decimal(5,2);not null;default:0"`
	TotalExpenses         float64                   `gorm:"column:total_expenses;type:decimal(12,2);not null;default:0"`
	TotalReceipts         int                       `gorm:"column:total_receipts;not null;default:0"`
	ReceiptsWithDocuments int                       `gorm:"column:receipts_with_documents;not null;default:0"`
	TotalActivities       int                       `gorm:"column:total_activities;not null;default:0"`
	TransparencyScore     float64                   `gorm:"column:transparency_score;type:decimal(5,2);not null;default:0"`
	TransparencyBreakdown TransparencyBreakdownJSON `gorm:"column:transparency_breakdown;type:jsonb"`
	AlertsCount           int                       `gorm:"column:alerts_count;not null;default:0"`
	AlertsResolved        int                       `gorm:"column:alerts_resolved;not null;default:0"`
	ReportPdfURL          *string                   `gorm:"column:report_pdf_url;type:text"`
	ReportHash            *string                   `gorm:"column:report_hash;type:varchar(64)"`
	ClosedAt              time.Time                 `gorm:"column:closed_at;not null"`
	CreatedAt             time.Time                 `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name for GORM
func (CampaignClosureReportModel) TableName() string {
	return "campaign_closure_reports"
}

// ToEntity converts a database model to a domain entity
func (m CampaignClosureReportModel) ToEntity() CampaignClosureReport {
	return CampaignClosureReport{
		ID:                    m.ID,
		CampaignID:            m.CampaignID,
		ClosureType:           ClosureType(m.ClosureType),
		ClosureReason:         m.ClosureReason,
		ClosedBy:              m.ClosedBy,
		TotalRaised:           m.TotalRaised,
		TotalDonors:           m.TotalDonors,
		TotalDonations:        m.TotalDonations,
		CampaignGoal:          m.CampaignGoal,
		GoalPercentage:        m.GoalPercentage,
		TotalExpenses:         m.TotalExpenses,
		TotalReceipts:         m.TotalReceipts,
		ReceiptsWithDocuments: m.ReceiptsWithDocuments,
		TotalActivities:       m.TotalActivities,
		TransparencyScore:     m.TransparencyScore,
		TransparencyBreakdown: TransparencyBreakdown(m.TransparencyBreakdown),
		AlertsCount:           m.AlertsCount,
		AlertsResolved:        m.AlertsResolved,
		ReportPdfURL:          m.ReportPdfURL,
		ReportHash:            m.ReportHash,
		ClosedAt:              m.ClosedAt,
		CreatedAt:             m.CreatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *CampaignClosureReportModel) FromEntity(entity CampaignClosureReport) {
	m.ID = entity.ID
	m.CampaignID = entity.CampaignID
	m.ClosureType = string(entity.ClosureType)
	m.ClosureReason = entity.ClosureReason
	m.ClosedBy = entity.ClosedBy
	m.TotalRaised = entity.TotalRaised
	m.TotalDonors = entity.TotalDonors
	m.TotalDonations = entity.TotalDonations
	m.CampaignGoal = entity.CampaignGoal
	m.GoalPercentage = entity.GoalPercentage
	m.TotalExpenses = entity.TotalExpenses
	m.TotalReceipts = entity.TotalReceipts
	m.ReceiptsWithDocuments = entity.ReceiptsWithDocuments
	m.TotalActivities = entity.TotalActivities
	m.TransparencyScore = entity.TransparencyScore
	m.TransparencyBreakdown = TransparencyBreakdownJSON(entity.TransparencyBreakdown)
	m.AlertsCount = entity.AlertsCount
	m.AlertsResolved = entity.AlertsResolved
	m.ReportPdfURL = entity.ReportPdfURL
	m.ReportHash = entity.ReportHash
	m.ClosedAt = entity.ClosedAt
	m.CreatedAt = entity.CreatedAt
}

// CampaignAlertModel represents the database table for alerts (placeholder)
type CampaignAlertModel struct {
	ID              uuid.UUID  `gorm:"primaryKey;column:id;type:uuid;default:uuid_generate_v4()"`
	CampaignID      uuid.UUID  `gorm:"column:campaign_id;type:uuid;not null;index"`
	AlertType       string     `gorm:"column:alert_type;type:varchar(50);not null"`
	Description     string     `gorm:"column:description;type:text;not null"`
	Status          string     `gorm:"column:status;type:varchar(20);not null;default:pending"`
	Severity        string     `gorm:"column:severity;type:varchar(20);not null;default:medium"`
	ReportedBy      *uuid.UUID `gorm:"column:reported_by;type:uuid"`
	ResolvedBy      *uuid.UUID `gorm:"column:resolved_by;type:uuid"`
	ResolutionNotes *string    `gorm:"column:resolution_notes;type:text"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
	ResolvedAt      *time.Time `gorm:"column:resolved_at"`
}

// TableName specifies the table name for GORM
func (CampaignAlertModel) TableName() string {
	return "campaign_alerts"
}
