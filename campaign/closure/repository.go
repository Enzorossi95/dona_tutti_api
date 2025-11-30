package closure

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Closure report operations
	GetClosureReport(ctx context.Context, campaignID uuid.UUID) (CampaignClosureReport, error)
	CreateClosureReport(ctx context.Context, report CampaignClosureReport) error
	UpdateReportPdfURL(ctx context.Context, campaignID uuid.UUID, pdfURL, hash string) error
	ExistsClosureReport(ctx context.Context, campaignID uuid.UUID) (bool, error)

	// Metrics queries for closure calculation
	GetDonationMetrics(ctx context.Context, campaignID uuid.UUID) (DonationMetrics, error)
	GetReceiptsMetrics(ctx context.Context, campaignID uuid.UUID) (ReceiptsMetrics, error)
	GetActivitiesMetrics(ctx context.Context, campaignID uuid.UUID) (ActivitiesMetrics, error)
	GetAlertsMetrics(ctx context.Context, campaignID uuid.UUID) (AlertsMetrics, error)

	// Data for PDF
	GetReceiptSummaries(ctx context.Context, campaignID uuid.UUID) ([]ReceiptSummary, error)
	GetActivitySummaries(ctx context.Context, campaignID uuid.UUID) ([]ActivitySummary, error)
}

// DonationMetrics holds donation statistics
type DonationMetrics struct {
	TotalRaised    float64
	TotalDonors    int
	TotalDonations int
}

// ReceiptsMetrics holds receipts statistics
type ReceiptsMetrics struct {
	TotalExpenses         float64
	TotalReceipts         int
	ReceiptsWithDocuments int
}

// ActivitiesMetrics holds activities statistics
type ActivitiesMetrics struct {
	TotalActivities              int
	AverageDaysBetweenActivities float64
}

// AlertsMetrics holds alerts statistics (placeholder)
type AlertsMetrics struct {
	AlertsCount    int
	AlertsResolved int
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetClosureReport(ctx context.Context, campaignID uuid.UUID) (CampaignClosureReport, error) {
	var model CampaignClosureReportModel
	if err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).First(&model).Error; err != nil {
		return CampaignClosureReport{}, err
	}
	return model.ToEntity(), nil
}

func (r *repository) CreateClosureReport(ctx context.Context, report CampaignClosureReport) error {
	var model CampaignClosureReportModel
	model.FromEntity(report)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *repository) UpdateReportPdfURL(ctx context.Context, campaignID uuid.UUID, pdfURL, hash string) error {
	return r.db.WithContext(ctx).
		Model(&CampaignClosureReportModel{}).
		Where("campaign_id = ?", campaignID).
		Updates(map[string]interface{}{
			"report_pdf_url": pdfURL,
			"report_hash":    hash,
		}).Error
}

func (r *repository) ExistsClosureReport(ctx context.Context, campaignID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&CampaignClosureReportModel{}).
		Where("campaign_id = ?", campaignID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) GetDonationMetrics(ctx context.Context, campaignID uuid.UUID) (DonationMetrics, error) {
	var result struct {
		TotalRaised    float64
		TotalDonors    int64
		TotalDonations int64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(amount), 0) as total_raised,
			COUNT(DISTINCT donor_id) as total_donors,
			COUNT(*) as total_donations
		FROM donations
		WHERE campaign_id = ? AND status = 'completed'
	`, campaignID).Scan(&result).Error

	if err != nil {
		return DonationMetrics{}, err
	}

	return DonationMetrics{
		TotalRaised:    result.TotalRaised,
		TotalDonors:    int(result.TotalDonors),
		TotalDonations: int(result.TotalDonations),
	}, nil
}

func (r *repository) GetReceiptsMetrics(ctx context.Context, campaignID uuid.UUID) (ReceiptsMetrics, error) {
	var result struct {
		TotalExpenses         float64
		TotalReceipts         int64
		ReceiptsWithDocuments int64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(total), 0) as total_expenses,
			COUNT(*) as total_receipts,
			COUNT(CASE WHEN document_url IS NOT NULL AND document_url != '' THEN 1 END) as receipts_with_documents
		FROM receipts
		WHERE campaign_id = ?
	`, campaignID).Scan(&result).Error

	if err != nil {
		return ReceiptsMetrics{}, err
	}

	return ReceiptsMetrics{
		TotalExpenses:         result.TotalExpenses,
		TotalReceipts:         int(result.TotalReceipts),
		ReceiptsWithDocuments: int(result.ReceiptsWithDocuments),
	}, nil
}

func (r *repository) GetActivitiesMetrics(ctx context.Context, campaignID uuid.UUID) (ActivitiesMetrics, error) {
	var result struct {
		TotalActivities int64
		AvgDays         float64
	}

	// Get total activities count
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) as total_activities
		FROM activities
		WHERE campaign_id = ?
	`, campaignID).Scan(&result.TotalActivities).Error

	if err != nil {
		return ActivitiesMetrics{}, err
	}

	// Calculate average days between activities
	if result.TotalActivities > 1 {
		err = r.db.WithContext(ctx).Raw(`
			WITH activity_dates AS (
				SELECT date, LAG(date) OVER (ORDER BY date) as prev_date
				FROM activities
				WHERE campaign_id = ?
				ORDER BY date
			)
			SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (date - prev_date)) / 86400), 30) as avg_days
			FROM activity_dates
			WHERE prev_date IS NOT NULL
		`, campaignID).Scan(&result.AvgDays).Error

		if err != nil {
			result.AvgDays = 30 // Default value
		}
	} else {
		result.AvgDays = 30 // Default if only one or no activities
	}

	return ActivitiesMetrics{
		TotalActivities:              int(result.TotalActivities),
		AverageDaysBetweenActivities: result.AvgDays,
	}, nil
}

func (r *repository) GetAlertsMetrics(ctx context.Context, campaignID uuid.UUID) (AlertsMetrics, error) {
	var result struct {
		AlertsCount    int64
		AlertsResolved int64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) as alerts_count,
			COUNT(CASE WHEN status = 'resolved' THEN 1 END) as alerts_resolved
		FROM campaign_alerts
		WHERE campaign_id = ?
	`, campaignID).Scan(&result).Error

	if err != nil {
		// Return zeros if alerts table doesn't exist or other error
		return AlertsMetrics{}, nil
	}

	return AlertsMetrics{
		AlertsCount:    int(result.AlertsCount),
		AlertsResolved: int(result.AlertsResolved),
	}, nil
}

func (r *repository) GetReceiptSummaries(ctx context.Context, campaignID uuid.UUID) ([]ReceiptSummary, error) {
	var results []struct {
		Provider    string
		Name        string
		Total       float64
		Date        string
		DocumentURL *string
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT provider, name, total, date, document_url
		FROM receipts
		WHERE campaign_id = ?
		ORDER BY date DESC
	`, campaignID).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	summaries := make([]ReceiptSummary, len(results))
	for i, r := range results {
		summaries[i] = ReceiptSummary{
			Provider:    r.Provider,
			Name:        r.Name,
			Total:       r.Total,
			HasDocument: r.DocumentURL != nil && *r.DocumentURL != "",
		}
	}

	return summaries, nil
}

func (r *repository) GetActivitySummaries(ctx context.Context, campaignID uuid.UUID) ([]ActivitySummary, error) {
	var results []struct {
		Title string
		Type  string
		Date  string
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT title, type, date
		FROM activities
		WHERE campaign_id = ?
		ORDER BY date DESC
	`, campaignID).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	summaries := make([]ActivitySummary, len(results))
	for i, r := range results {
		summaries[i] = ActivitySummary{
			Title: r.Title,
			Type:  r.Type,
		}
	}

	return summaries, nil
}
