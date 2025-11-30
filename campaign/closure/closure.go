package closure

import (
	"time"

	"github.com/google/uuid"
)

// ClosureType represents how the campaign was closed
type ClosureType string

const (
	ClosureTypeGoalReached ClosureType = "goal_reached"
	ClosureTypeEndDate     ClosureType = "end_date"
	ClosureTypeManual      ClosureType = "manual"
)

// TransparencyBreakdown details the score calculation
type TransparencyBreakdown struct {
	DocumentationScore   float64 `json:"documentation_score"`    // 0-30 pts
	ActivityScore        float64 `json:"activity_score"`         // 0-25 pts
	GoalProgressScore    float64 `json:"goal_progress_score"`    // 0-20 pts
	TimelinessScore      float64 `json:"timeliness_score"`       // 0-15 pts
	AlertsDeductionScore float64 `json:"alerts_deduction_score"` // 0 to -10 pts
	BonusScore           float64 `json:"bonus_score"`            // 0-10 pts
}

// Total calculates the total transparency score
func (b TransparencyBreakdown) Total() float64 {
	total := b.DocumentationScore + b.ActivityScore + b.GoalProgressScore +
		b.TimelinessScore + b.AlertsDeductionScore + b.BonusScore
	if total < 0 {
		return 0
	}
	if total > 100 {
		return 100
	}
	return total
}

// CampaignClosureReport represents the full closure report
type CampaignClosureReport struct {
	ID                    uuid.UUID             `json:"id"`
	CampaignID            uuid.UUID             `json:"campaign_id"`
	ClosureType           ClosureType           `json:"closure_type"`
	ClosureReason         *string               `json:"closure_reason,omitempty"`
	ClosedBy              *uuid.UUID            `json:"closed_by,omitempty"`
	TotalRaised           float64               `json:"total_raised"`
	TotalDonors           int                   `json:"total_donors"`
	TotalDonations        int                   `json:"total_donations"`
	CampaignGoal          float64               `json:"campaign_goal"`
	GoalPercentage        float64               `json:"goal_percentage"`
	TotalExpenses         float64               `json:"total_expenses"`
	TotalReceipts         int                   `json:"total_receipts"`
	ReceiptsWithDocuments int                   `json:"receipts_with_documents"`
	TotalActivities       int                   `json:"total_activities"`
	TransparencyScore     float64               `json:"transparency_score"`
	TransparencyBreakdown TransparencyBreakdown `json:"transparency_breakdown"`
	AlertsCount           int                   `json:"alerts_count"`
	AlertsResolved        int                   `json:"alerts_resolved"`
	ReportPdfURL          *string               `json:"report_pdf_url,omitempty"`
	ReportHash            *string               `json:"report_hash,omitempty"`
	ClosedAt              time.Time             `json:"closed_at"`
	CreatedAt             time.Time             `json:"created_at"`
}

// PublicAuditReport is the public version for donors
type PublicAuditReport struct {
	CampaignID        uuid.UUID `json:"campaign_id"`
	CampaignTitle     string    `json:"campaign_title"`
	OrganizerName     string    `json:"organizer_name"`
	ClosedAt          time.Time `json:"closed_at"`
	TotalRaised       float64   `json:"total_raised"`
	CampaignGoal      float64   `json:"campaign_goal"`
	GoalPercentage    float64   `json:"goal_percentage"`
	TotalDonors       int       `json:"total_donors"`
	TotalExpenses     float64   `json:"total_expenses"`
	TransparencyScore float64   `json:"transparency_score"`
	ReportPdfURL      *string   `json:"report_pdf_url,omitempty"`
}

// CloseCampaignRequest for manual closure
type CloseCampaignRequest struct {
	Reason string `json:"reason" validate:"required,min=10"`
}

// ClosureMetrics holds the data needed to calculate the closure report
type ClosureMetrics struct {
	// Campaign data
	CampaignGoal   float64
	CampaignStart  time.Time
	CampaignEnd    time.Time
	HasContract    bool
	CampaignTitle  string
	OrganizerName  string
	OrganizerID    uuid.UUID

	// Donations
	TotalRaised    float64
	TotalDonors    int
	TotalDonations int

	// Receipts
	TotalExpenses         float64
	TotalReceipts         int
	ReceiptsWithDocuments int

	// Activities
	TotalActivities              int
	AverageDaysBetweenActivities float64

	// Alerts (placeholder)
	AlertsCount    int
	AlertsResolved int
}

// AuditReportData contains all data needed to generate the PDF
type AuditReportData struct {
	CampaignID      uuid.UUID
	CampaignTitle   string
	CampaignGoal    float64
	OrganizerName   string
	StartDate       time.Time
	EndDate         time.Time
	ClosedAt        time.Time
	ClosureType     ClosureType
	ClosureReason   *string

	// Financial
	TotalRaised    float64
	GoalPercentage float64
	TotalDonors    int
	TotalDonations int

	// Expenses
	TotalExpenses         float64
	TotalReceipts         int
	ReceiptsWithDocuments int
	Receipts              []ReceiptSummary

	// Activities
	TotalActivities int
	Activities      []ActivitySummary

	// Transparency
	TransparencyScore     float64
	TransparencyBreakdown TransparencyBreakdown
}

// ReceiptSummary for PDF display
type ReceiptSummary struct {
	Provider    string    `json:"provider"`
	Name        string    `json:"name"`
	Total       float64   `json:"total"`
	Date        time.Time `json:"date"`
	HasDocument bool      `json:"has_document"`
}

// ActivitySummary for PDF display
type ActivitySummary struct {
	Title string    `json:"title"`
	Type  string    `json:"type"`
	Date  time.Time `json:"date"`
}
