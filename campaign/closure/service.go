package closure

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"time"

	"dona_tutti_api/s3client"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Service defines the interface for campaign closure business logic
type Service interface {
	CloseCampaign(ctx context.Context, campaignID uuid.UUID, closureType ClosureType, reason *string, closedBy *uuid.UUID) (*CampaignClosureReport, error)
	GetClosureReport(ctx context.Context, campaignID uuid.UUID) (*CampaignClosureReport, error)
	GetPublicAuditReport(ctx context.Context, campaignID uuid.UUID) (*PublicAuditReport, error)
	HasClosureReport(ctx context.Context, campaignID uuid.UUID) (bool, error)
}

// CampaignInfo represents minimal campaign information needed for closure
type CampaignInfo struct {
	ID          uuid.UUID
	Title       string
	Goal        float64
	OrganizerID uuid.UUID
	Status      string
	StartDate   time.Time
	EndDate     time.Time
}

// CampaignServiceInterface defines the interface for campaign operations
type CampaignServiceInterface interface {
	GetCampaignForClosure(ctx context.Context, id uuid.UUID) (CampaignInfo, error)
	UpdateStatus(ctx context.Context, campaignID uuid.UUID, status string) error
}

// OrganizerServiceInterface defines the interface for organizer operations
type OrganizerServiceInterface interface {
	GetOrganizerName(ctx context.Context, organizerID uuid.UUID) (string, error)
}

// ContractServiceInterface defines the interface for contract operations
type ContractServiceInterface interface {
	HasContract(ctx context.Context, campaignID uuid.UUID) (bool, error)
}

// PDFGenerator defines the interface for generating PDF audit reports
type PDFGenerator interface {
	Generate(data AuditReportData) ([]byte, string, error) // returns PDF bytes, hash, error
}

type service struct {
	repo             Repository
	pdfGenerator     PDFGenerator
	s3Client         *s3client.Client
	campaignService  CampaignServiceInterface
	organizerService OrganizerServiceInterface
	contractService  ContractServiceInterface
}

// NewService creates a new instance of the closure service
func NewService(
	repo Repository,
	pdfGenerator PDFGenerator,
	s3Client *s3client.Client,
	campaignService CampaignServiceInterface,
	organizerService OrganizerServiceInterface,
	contractService ContractServiceInterface,
) Service {
	return &service{
		repo:             repo,
		pdfGenerator:     pdfGenerator,
		s3Client:         s3Client,
		campaignService:  campaignService,
		organizerService: organizerService,
		contractService:  contractService,
	}
}

// CloseCampaign closes a campaign and generates the audit report
func (s *service) CloseCampaign(ctx context.Context, campaignID uuid.UUID, closureType ClosureType, reason *string, closedBy *uuid.UUID) (*CampaignClosureReport, error) {
	// 1. Check if closure report already exists
	exists, err := s.repo.ExistsClosureReport(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to check closure report: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("campaign %s already has a closure report", campaignID)
	}

	// 2. Get campaign info
	campaignInfo, err := s.campaignService.GetCampaignForClosure(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// 3. Validate campaign can be closed
	if campaignInfo.Status != "active" && campaignInfo.Status != "paused" {
		return nil, fmt.Errorf("campaign must be in active or paused status to close, current status: %s", campaignInfo.Status)
	}

	// 4. Validate manual closure has reason
	if closureType == ClosureTypeManual && (reason == nil || len(*reason) < 10) {
		return nil, fmt.Errorf("manual closure requires a reason with at least 10 characters")
	}

	// 5. Get organizer name
	organizerName, err := s.organizerService.GetOrganizerName(ctx, campaignInfo.OrganizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizer: %w", err)
	}

	// 6. Gather all metrics
	donationMetrics, err := s.repo.GetDonationMetrics(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get donation metrics: %w", err)
	}

	receiptsMetrics, err := s.repo.GetReceiptsMetrics(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipts metrics: %w", err)
	}

	activitiesMetrics, err := s.repo.GetActivitiesMetrics(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities metrics: %w", err)
	}

	alertsMetrics, _ := s.repo.GetAlertsMetrics(ctx, campaignID) // Ignore error for placeholder

	// 7. Check if campaign has contract
	hasContract, _ := s.contractService.HasContract(ctx, campaignID)

	// 8. Calculate transparency score
	closureMetrics := ClosureMetrics{
		CampaignGoal:                 campaignInfo.Goal,
		CampaignStart:                campaignInfo.StartDate,
		CampaignEnd:                  campaignInfo.EndDate,
		HasContract:                  hasContract,
		CampaignTitle:                campaignInfo.Title,
		OrganizerName:                organizerName,
		OrganizerID:                  campaignInfo.OrganizerID,
		TotalRaised:                  donationMetrics.TotalRaised,
		TotalDonors:                  donationMetrics.TotalDonors,
		TotalDonations:               donationMetrics.TotalDonations,
		TotalExpenses:                receiptsMetrics.TotalExpenses,
		TotalReceipts:                receiptsMetrics.TotalReceipts,
		ReceiptsWithDocuments:        receiptsMetrics.ReceiptsWithDocuments,
		TotalActivities:              activitiesMetrics.TotalActivities,
		AverageDaysBetweenActivities: activitiesMetrics.AverageDaysBetweenActivities,
		AlertsCount:                  alertsMetrics.AlertsCount,
		AlertsResolved:               alertsMetrics.AlertsResolved,
	}

	breakdown := s.calculateTransparencyScore(closureMetrics, closureType == ClosureTypeManual)
	transparencyScore := breakdown.Total()

	// 9. Calculate goal percentage
	goalPercentage := 0.0
	if campaignInfo.Goal > 0 {
		goalPercentage = (donationMetrics.TotalRaised / campaignInfo.Goal) * 100
		if goalPercentage > 100 {
			goalPercentage = 100
		}
	}

	// 10. Create closure report
	closedAt := time.Now()
	report := CampaignClosureReport{
		ID:                    uuid.New(),
		CampaignID:            campaignID,
		ClosureType:           closureType,
		ClosureReason:         reason,
		ClosedBy:              closedBy,
		TotalRaised:           donationMetrics.TotalRaised,
		TotalDonors:           donationMetrics.TotalDonors,
		TotalDonations:        donationMetrics.TotalDonations,
		CampaignGoal:          campaignInfo.Goal,
		GoalPercentage:        goalPercentage,
		TotalExpenses:         receiptsMetrics.TotalExpenses,
		TotalReceipts:         receiptsMetrics.TotalReceipts,
		ReceiptsWithDocuments: receiptsMetrics.ReceiptsWithDocuments,
		TotalActivities:       activitiesMetrics.TotalActivities,
		TransparencyScore:     transparencyScore,
		TransparencyBreakdown: breakdown,
		AlertsCount:           alertsMetrics.AlertsCount,
		AlertsResolved:        alertsMetrics.AlertsResolved,
		ClosedAt:              closedAt,
		CreatedAt:             closedAt,
	}

	// 11. Save closure report
	if err := s.repo.CreateClosureReport(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create closure report: %w", err)
	}

	// 12. Update campaign status to completed
	if err := s.campaignService.UpdateStatus(ctx, campaignID, "completed"); err != nil {
		return nil, fmt.Errorf("failed to update campaign status: %w", err)
	}

	// 13. Generate PDF asynchronously
	go s.generateAndUploadPDF(context.Background(), campaignID, campaignInfo, organizerName, report, closureMetrics)

	return &report, nil
}

// generateAndUploadPDF generates the audit PDF and uploads it to S3
func (s *service) generateAndUploadPDF(ctx context.Context, campaignID uuid.UUID, campaignInfo CampaignInfo, organizerName string, report CampaignClosureReport, metrics ClosureMetrics) {
	// Get receipt and activity summaries for PDF
	receiptSummaries, _ := s.repo.GetReceiptSummaries(ctx, campaignID)
	activitySummaries, _ := s.repo.GetActivitySummaries(ctx, campaignID)

	// Build audit report data
	data := AuditReportData{
		CampaignID:            campaignID,
		CampaignTitle:         campaignInfo.Title,
		CampaignGoal:          campaignInfo.Goal,
		OrganizerName:         organizerName,
		StartDate:             campaignInfo.StartDate,
		EndDate:               campaignInfo.EndDate,
		ClosedAt:              report.ClosedAt,
		ClosureType:           report.ClosureType,
		ClosureReason:         report.ClosureReason,
		TotalRaised:           report.TotalRaised,
		GoalPercentage:        report.GoalPercentage,
		TotalDonors:           report.TotalDonors,
		TotalDonations:        report.TotalDonations,
		TotalExpenses:         report.TotalExpenses,
		TotalReceipts:         report.TotalReceipts,
		ReceiptsWithDocuments: report.ReceiptsWithDocuments,
		Receipts:              receiptSummaries,
		TotalActivities:       report.TotalActivities,
		Activities:            activitySummaries,
		TransparencyScore:     report.TransparencyScore,
		TransparencyBreakdown: report.TransparencyBreakdown,
	}

	// Generate PDF
	pdfBytes, hash, err := s.pdfGenerator.Generate(data)
	if err != nil {
		fmt.Printf("failed to generate audit PDF for campaign %s: %v\n", campaignID, err)
		return
	}

	// Generate S3 key
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("audits/%s/audit-report-%d.pdf", campaignID, timestamp)

	// Upload to S3
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.s3Client.GetBucketName()),
		Key:         aws.String(key),
		Body:        bytes.NewReader(pdfBytes),
		ContentType: aws.String("application/pdf"),
	}

	_, err = s.s3Client.GetS3Client().PutObject(ctx, uploadInput)
	if err != nil {
		fmt.Printf("failed to upload audit PDF for campaign %s: %v\n", campaignID, err)
		return
	}

	// Generate public URL
	var url string
	if s.s3Client.GetEndpoint() != "" {
		url = fmt.Sprintf("%s/%s/%s", s.s3Client.GetEndpoint(), s.s3Client.GetBucketName(), key)
	} else {
		url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.s3Client.GetBucketName(), key)
	}

	// Update report with PDF URL
	if err := s.repo.UpdateReportPdfURL(ctx, campaignID, url, hash); err != nil {
		fmt.Printf("failed to update audit PDF URL for campaign %s: %v\n", campaignID, err)
	}
}

// calculateTransparencyScore calculates the transparency score based on metrics
func (s *service) calculateTransparencyScore(metrics ClosureMetrics, closedBeforeEndDate bool) TransparencyBreakdown {
	breakdown := TransparencyBreakdown{}

	// 1. DOCUMENTATION (30 pts max)
	// Percentage of receipts with attached document
	// If no receipts exist, score remains 0 (default value)
	if metrics.TotalReceipts > 0 {
		docPercentage := float64(metrics.ReceiptsWithDocuments) / float64(metrics.TotalReceipts)
		breakdown.DocumentationScore = docPercentage * 30.0
	}

	// 2. ACTIVITIES (25 pts max)
	// Minimum 1 activity per month of active campaign
	campaignDurationMonths := math.Ceil(metrics.CampaignEnd.Sub(metrics.CampaignStart).Hours() / (24 * 30))
	if campaignDurationMonths < 1 {
		campaignDurationMonths = 1
	}
	expectedActivities := int(campaignDurationMonths)
	if metrics.TotalActivities > 0 {
		activityRatio := float64(metrics.TotalActivities) / float64(expectedActivities)
		if activityRatio > 1 {
			activityRatio = 1 + (activityRatio-1)*0.2 // Bonus for more activities
		}
		if activityRatio > 1.5 {
			activityRatio = 1.5
		}
		breakdown.ActivityScore = activityRatio * 25.0 / 1.5
		if breakdown.ActivityScore > 25 {
			breakdown.ActivityScore = 25
		}
	}

	// 3. GOAL PROGRESS (20 pts max)
	// Rewards reaching or getting close to goal
	if metrics.CampaignGoal > 0 {
		goalRatio := metrics.TotalRaised / metrics.CampaignGoal
		if goalRatio >= 1 {
			breakdown.GoalProgressScore = 20.0
		} else if goalRatio >= 0.75 {
			breakdown.GoalProgressScore = 15.0 + (goalRatio-0.75)*20.0
		} else if goalRatio >= 0.5 {
			breakdown.GoalProgressScore = 10.0 + (goalRatio-0.5)*20.0
		} else {
			breakdown.GoalProgressScore = goalRatio * 20.0
		}
	}

	// 4. TIMELINESS (15 pts max)
	// Based on frequency of activities
	avgDays := metrics.AverageDaysBetweenActivities
	if avgDays <= 7 {
		breakdown.TimelinessScore = 15.0
	} else if avgDays <= 14 {
		breakdown.TimelinessScore = 12.0
	} else if avgDays <= 30 {
		breakdown.TimelinessScore = 8.0
	} else {
		breakdown.TimelinessScore = 5.0
	}

	// 5. ALERTS DEDUCTION (0 to -10 pts) - PLACEHOLDER
	if metrics.AlertsCount > 0 {
		unresolvedAlerts := metrics.AlertsCount - metrics.AlertsResolved
		breakdown.AlertsDeductionScore = float64(unresolvedAlerts) * -2.0
		if breakdown.AlertsDeductionScore < -10 {
			breakdown.AlertsDeductionScore = -10
		}
	}

	// 6. BONUS (10 pts max)
	bonus := 0.0
	if metrics.HasContract {
		bonus += 3.0
	}
	if metrics.TotalDonors >= 10 {
		bonus += 2.0
	}
	if metrics.TotalRaised > 0 && metrics.TotalExpenses > 0 {
		expenseRatio := metrics.TotalExpenses / metrics.TotalRaised
		if expenseRatio >= 0.8 {
			bonus += 3.0
		}
	}
	if closedBeforeEndDate {
		bonus += 2.0
	}
	breakdown.BonusScore = bonus

	return breakdown
}

// GetClosureReport retrieves the closure report for a campaign
func (s *service) GetClosureReport(ctx context.Context, campaignID uuid.UUID) (*CampaignClosureReport, error) {
	report, err := s.repo.GetClosureReport(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("closure report not found: %w", err)
	}
	return &report, nil
}

// GetPublicAuditReport retrieves the public audit report for donors
func (s *service) GetPublicAuditReport(ctx context.Context, campaignID uuid.UUID) (*PublicAuditReport, error) {
	// Get closure report
	report, err := s.repo.GetClosureReport(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("audit report not found: %w", err)
	}

	// Get campaign info
	campaignInfo, err := s.campaignService.GetCampaignForClosure(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Get organizer name
	organizerName, err := s.organizerService.GetOrganizerName(ctx, campaignInfo.OrganizerID)
	if err != nil {
		organizerName = "Unknown"
	}

	return &PublicAuditReport{
		CampaignID:        campaignID,
		CampaignTitle:     campaignInfo.Title,
		OrganizerName:     organizerName,
		ClosedAt:          report.ClosedAt,
		TotalRaised:       report.TotalRaised,
		CampaignGoal:      report.CampaignGoal,
		GoalPercentage:    report.GoalPercentage,
		TotalDonors:       report.TotalDonors,
		TotalExpenses:     report.TotalExpenses,
		TransparencyScore: report.TransparencyScore,
		ReportPdfURL:      report.ReportPdfURL,
	}, nil
}

// HasClosureReport checks if a closure report exists for a campaign
func (s *service) HasClosureReport(ctx context.Context, campaignID uuid.UUID) (bool, error) {
	return s.repo.ExistsClosureReport(ctx, campaignID)
}
