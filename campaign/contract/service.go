package contract

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"dona_tutti_api/s3client"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Service defines the interface for campaign contract business logic
type Service interface {
	GenerateContract(ctx context.Context, campaignID uuid.UUID) (string, error)
	GetContract(ctx context.Context, campaignID uuid.UUID) (CampaignContract, error)
	AcceptContract(ctx context.Context, req AcceptContractRequest) error
	GetContractProof(ctx context.Context, campaignID uuid.UUID) (ContractProof, error)
	HasContract(ctx context.Context, campaignID uuid.UUID) (bool, error)
}

// PDFGenerator defines the interface for generating PDF contracts
type PDFGenerator interface {
	Generate(data ContractData) ([]byte, string, error) // returns PDF bytes, hash, error
}

// CampaignInfo represents minimal campaign information needed for contracts
type CampaignInfo struct {
	ID          uuid.UUID
	Title       string
	Goal        float64
	OrganizerID uuid.UUID
	Status      string
}

// OrganizerInfo represents minimal organizer information needed for contracts
type OrganizerInfo struct {
	ID      uuid.UUID
	Name    string
	Email   string
	Phone   string
	Address string
}

// CampaignService defines the interface for campaign operations
type CampaignService interface {
	GetCampaignInfo(ctx context.Context, id uuid.UUID) (CampaignInfo, error)
	UpdateStatus(ctx context.Context, campaignID uuid.UUID, status string) error
	GetCampaignTitle(ctx context.Context, campaignID uuid.UUID) (string, error)
}

// OrganizerService defines the interface for organizer operations
type OrganizerService interface {
	GetOrganizerInfo(ctx context.Context, id uuid.UUID) (OrganizerInfo, error)
	GetOrganizerName(ctx context.Context, organizerID uuid.UUID) (string, error)
}

type service struct {
	repo             Repository
	pdfGenerator     PDFGenerator
	s3Client         *s3client.Client
	campaignService  CampaignService
	organizerService OrganizerService
}

// NewService creates a new instance of the contract service
func NewService(
	repo Repository,
	pdfGenerator PDFGenerator,
	s3Client *s3client.Client,
	campaignService CampaignService,
	organizerService OrganizerService,
) Service {
	return &service{
		repo:             repo,
		pdfGenerator:     pdfGenerator,
		s3Client:         s3Client,
		campaignService:  campaignService,
		organizerService: organizerService,
	}
}

// GenerateContract generates a contract PDF and uploads it to S3
func (s *service) GenerateContract(ctx context.Context, campaignID uuid.UUID) (string, error) {
	// 1. Check if contract already exists
	exists, err := s.repo.ExistsByCampaignID(ctx, campaignID)
	if err != nil {
		return "", fmt.Errorf("failed to check contract existence: %w", err)
	}
	if exists {
		return "", fmt.Errorf("contract already exists for campaign %s", campaignID)
	}

	// 2. Fetch campaign info from database
	campaignInfo, err := s.campaignService.GetCampaignInfo(ctx, campaignID)
	if err != nil {
		return "", fmt.Errorf("campaign not found: %w", err)
	}

	// 3. Validate campaign status (must be draft to generate contract)
	if campaignInfo.Status != "draft" {
		return "", fmt.Errorf("contract can only be generated for campaigns in draft status, current status: %s", campaignInfo.Status)
	}

	// 4. Validate campaign data
	if campaignInfo.Title == "" {
		return "", fmt.Errorf("campaign must have a title")
	}
	if campaignInfo.Goal <= 0 {
		return "", fmt.Errorf("campaign must have a valid goal")
	}
	if campaignInfo.OrganizerID == uuid.Nil {
		return "", fmt.Errorf("campaign must have an associated organizer")
	}

	// 5. Fetch organizer info from database
	organizerInfo, err := s.organizerService.GetOrganizerInfo(ctx, campaignInfo.OrganizerID)
	if err != nil {
		return "", fmt.Errorf("organizer not found: %w", err)
	}

	// 6. Validate required organizer data
	if organizerInfo.Name == "" {
		return "", fmt.Errorf("organizer must have a name")
	}
	if organizerInfo.Email == "" {
		return "", fmt.Errorf("organizer must have an email")
	}
	if organizerInfo.Phone == "" {
		return "", fmt.Errorf("organizer must have a phone")
	}
	if organizerInfo.Address == "" {
		return "", fmt.Errorf("organizer must have an address")
	}

	// 7. Build contract data from database info
	data := ContractData{
		CampaignID:       campaignInfo.ID,
		CampaignTitle:    campaignInfo.Title,
		CampaignGoal:     campaignInfo.Goal,
		OrganizerID:      organizerInfo.ID,
		OrganizerName:    organizerInfo.Name,
		OrganizerEmail:   organizerInfo.Email,
		OrganizerPhone:   organizerInfo.Phone,
		OrganizerAddress: organizerInfo.Address,
		GeneratedAt:      time.Now(),
	}

	// 8. Generate PDF
	pdfBytes, hash, err := s.pdfGenerator.Generate(data)

	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// 9. Generate S3 key
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("contracts/%s/contract-%d.pdf", campaignID, timestamp)

	// 10. Upload to S3
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.s3Client.GetBucketName()),
		Key:         aws.String(key),
		Body:        bytes.NewReader(pdfBytes),
		ContentType: aws.String("application/pdf"),
	}

	_, err = s.s3Client.GetS3Client().PutObject(ctx, uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload contract to S3: %w", err)
	}

	// 11. Generate public URL
	var url string
	if s.s3Client.GetEndpoint() != "" {
		// LocalStack URL
		url = fmt.Sprintf("%s/%s/%s", s.s3Client.GetEndpoint(), s.s3Client.GetBucketName(), key)
	} else {
		// AWS URL
		url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.s3Client.GetBucketName(), key)
	}

	// 12. Store contract metadata (without acceptance yet)
	contract := CampaignContract{
		ID:             uuid.New(),
		CampaignID:     campaignID,
		OrganizerID:    organizerInfo.ID,
		ContractPdfURL: url,
		ContractHash:   hash,
		AcceptedAt:     time.Time{}, // Not accepted yet
		Acceptance: AcceptanceMetadata{
			IP:        "",
			UserAgent: "",
		},
		CreatedAt: time.Now(),
	}

	// 13. Save the contract metadata
	if err := s.repo.Create(ctx, contract); err != nil {
		return "", fmt.Errorf("failed to save contract metadata: %w", err)
	}

	// 14. Update campaign status to pending_approval
	if err := s.campaignService.UpdateStatus(ctx, campaignID, "pending_approval"); err != nil {
		return "", fmt.Errorf("failed to update campaign status: %w", err)
	}

	return url, nil
}

// GetContract retrieves a contract by campaign ID
func (s *service) GetContract(ctx context.Context, campaignID uuid.UUID) (CampaignContract, error) {
	contract, err := s.repo.GetByCampaignID(ctx, campaignID)
	if err != nil {
		return CampaignContract{}, fmt.Errorf("contract not found: %w", err)
	}
	return contract, nil
}

// AcceptContract records the acceptance of a contract
func (s *service) AcceptContract(ctx context.Context, req AcceptContractRequest) error {
	// Validate request
	if req.IP == "" {
		return fmt.Errorf("acceptance IP is required")
	}

	// Get existing contract (must be generated first)
	contract, err := s.repo.GetByCampaignID(ctx, req.CampaignID)
	if err != nil {
		return fmt.Errorf("contract not found - must generate contract first: %w", err)
	}

	// Verify campaign status (must be pending_approval to accept contract)
	campaignInfo, err := s.campaignService.GetCampaignInfo(ctx, req.CampaignID)
	if err != nil {
		return fmt.Errorf("campaign not found: %w", err)
	}
	if campaignInfo.Status != "pending_approval" {
		return fmt.Errorf("contract can only be accepted for campaigns in pending_approval status, current status: %s", campaignInfo.Status)
	}

	// Check if already accepted
	if !contract.AcceptedAt.IsZero() {
		return fmt.Errorf("contract already accepted for campaign %s", req.CampaignID)
	}

	// Update contract with acceptance metadata
	contract.AcceptedAt = time.Now()
	contract.Acceptance = AcceptanceMetadata{
		IP:        req.IP,
		UserAgent: req.UserAgent,
	}

	// Save the updated contract
	if err := s.repo.Update(ctx, contract); err != nil {
		return fmt.Errorf("failed to update contract with acceptance: %w", err)
	}

	// Update campaign status to active (published)
	if err := s.campaignService.UpdateStatus(ctx, req.CampaignID, "active"); err != nil {
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	return nil
}

// GetContractProof retrieves the contract proof for admin review
func (s *service) GetContractProof(ctx context.Context, campaignID uuid.UUID) (ContractProof, error) {
	// Get contract
	contract, err := s.repo.GetByCampaignID(ctx, campaignID)
	if err != nil {
		return ContractProof{}, fmt.Errorf("contract not found: %w", err)
	}

	// Get campaign title
	campaignTitle, err := s.campaignService.GetCampaignTitle(ctx, campaignID)
	if err != nil {
		return ContractProof{}, fmt.Errorf("failed to get campaign title: %w", err)
	}

	// Get organizer name
	organizerName, err := s.organizerService.GetOrganizerName(ctx, contract.OrganizerID)
	if err != nil {
		return ContractProof{}, fmt.Errorf("failed to get organizer name: %w", err)
	}

	return ContractProof{
		Contract:      contract,
		CampaignTitle: campaignTitle,
		OrganizerName: organizerName,
	}, nil
}

// HasContract checks if a contract exists for a campaign
func (s *service) HasContract(ctx context.Context, campaignID uuid.UUID) (bool, error) {
	return s.repo.ExistsByCampaignID(ctx, campaignID)
}
