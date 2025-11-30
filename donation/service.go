package donation

import (
	"context"
	"dona_tutti_api/donation/receipt"
	"dona_tutti_api/donor"
	"dona_tutti_api/s3client"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetDonation(ctx context.Context, id uuid.UUID) (Donation, error)
	CreateDonation(ctx context.Context, donation Donation) (uuid.UUID, error)
	CreateDonationWithRequest(ctx context.Context, campaignID uuid.UUID, req CreateDonationRequest) (uuid.UUID, error)
	UpdateDonation(ctx context.Context, donation Donation) error
	UpdateDonationStatus(ctx context.Context, id uuid.UUID, status DonationStatus) error
	ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error)
}

// CampaignService defines minimal campaign operations needed by donation service
type CampaignService interface {
	GetCampaignTitle(ctx context.Context, campaignID uuid.UUID) (string, error)
}

type service struct {
	repo            DonationRepository
	donorService    donor.Service
	s3Client        *s3client.Client
	campaignService CampaignService
	pdfGenerator    receipt.PDFGenerator
}

func NewService(repo DonationRepository, donorService donor.Service, s3Client *s3client.Client, campaignService CampaignService) Service {
	return &service{
		repo:            repo,
		donorService:    donorService,
		s3Client:        s3Client,
		campaignService: campaignService,
		pdfGenerator:    receipt.NewPDFGenerator(),
	}
}

func (s *service) GetDonation(ctx context.Context, id uuid.UUID) (Donation, error) {
	return s.repo.GetDonation(ctx, id)
}

func (s *service) CreateDonation(ctx context.Context, donation Donation) (uuid.UUID, error) {
	donation.ID = uuid.New()
	donation.Date = time.Now()

	if donation.Status == "" {
		donation.Status = DonationStatusPending
	}

	if err := s.repo.CreateDonation(ctx, donation); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create donation: %w", err)
	}

	// Note: Receipt generation is now handled by UpdateDonationStatus when status changes to completed
	// This allows for better control over when receipts are generated

	return donation.ID, nil
}

func (s *service) UpdateDonation(ctx context.Context, donation Donation) error {
	return s.repo.UpdateDonation(ctx, donation)
}

func (s *service) UpdateDonationStatus(ctx context.Context, id uuid.UUID, status DonationStatus) error {
	// Validate status
	if !IsValidStatus(status) {
		return fmt.Errorf("invalid donation status: %s", status)
	}

	// Get current donation
	currentDonation, err := s.repo.GetDonation(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get donation: %w", err)
	}

	// Business rule: Don't allow changing from completed to pending
	if currentDonation.Status == DonationStatusCompleted && status == DonationStatusPending {
		return fmt.Errorf("cannot change donation status from completed to pending")
	}

	// Update status in database
	currentDonation.Status = status
	if err := s.repo.UpdateDonation(ctx, currentDonation); err != nil {
		return fmt.Errorf("failed to update donation status: %w", err)
	}

	// If status changed to completed and no receipt exists, generate receipt
	if status == DonationStatusCompleted && currentDonation.ReceiptURL == nil {
		log.Printf("🎫 Donation %s marked as completed, generating receipt...", id.String())
		go s.generateAndUploadReceipt(context.Background(), currentDonation)
	}

	return nil
}

func (s *service) ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error) {
	return s.repo.ListDonationsByCampaign(ctx, campaignID)
}

func (s *service) GetOrCreateDonor(ctx context.Context, donorInfo DonorInfo) (uuid.UUID, error) {
	if donorInfo.Email != nil && *donorInfo.Email != "" {
		existingDonor, err := s.donorService.FindDonorByEmail(ctx, *donorInfo.Email)
		if err == nil {
			return existingDonor.ID, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, fmt.Errorf("error searching donor by email: %w", err)
		}
	}

	if donorInfo.Phone != nil && *donorInfo.Phone != "" {
		existingDonor, err := s.donorService.FindDonorByPhone(ctx, *donorInfo.Phone)
		if err == nil {
			return existingDonor.ID, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, fmt.Errorf("error searching donor by phone: %w", err)
		}
	}

	newDonor := donor.Donor{
		FirstName: donorInfo.Name,
		LastName:  donorInfo.LastName,
		Email:     "",
		Phone:     "",
	}
	if donorInfo.Email != nil {
		newDonor.Email = *donorInfo.Email
	}
	if donorInfo.Phone != nil {
		newDonor.Phone = *donorInfo.Phone
	}

	donorID, err := s.donorService.CreateDonor(ctx, newDonor)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create new donor: %w", err)
	}

	return donorID, nil
}

func (s *service) CreateDonationWithRequest(ctx context.Context, campaignID uuid.UUID, req CreateDonationRequest) (uuid.UUID, error) {
	var donorID uuid.UUID
	var err error

	if req.DonorID != nil {
		if req.IsAnonymous {
			return uuid.Nil, fmt.Errorf("donation with donor_id cannot be anonymous")
		}
		donorID = *req.DonorID
	} else if req.Donor != nil {
		if req.IsAnonymous {
			return uuid.Nil, fmt.Errorf("donation with donor information cannot be anonymous")
		}
		donorID, err = s.GetOrCreateDonor(ctx, *req.Donor)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to get or create donor: %w", err)
		}
	} else {
		if !req.IsAnonymous {
			return uuid.Nil, fmt.Errorf("anonymous donation must have is_anonymous set to true")
		}
		// Para donaciones anónimas, crear un donor temporal con información mínima
		anonymousDonor := donor.Donor{
			FirstName: "Anonymous",
			LastName:  "Donor",
			Email:     "",
			Phone:     "",
		}
		donorID, err = s.donorService.CreateDonor(ctx, anonymousDonor)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to create anonymous donor: %w", err)
		}
	}

	donation := Donation{
		CampaignID:      campaignID,
		Amount:          req.Amount,
		DonorID:         donorID,
		Message:         req.Message,
		IsAnonymous:     req.IsAnonymous,
		PaymentMethodID: req.PaymentMethodID,
		Status:          DonationStatusPending,
	}

	return s.CreateDonation(ctx, donation)
}

// generateAndUploadReceipt generates a PDF receipt and uploads it to S3
func (s *service) generateAndUploadReceipt(ctx context.Context, donation Donation) {
	// Skip if S3 client is not available
	if s.s3Client == nil {
		log.Printf("⚠️  Cannot generate receipt for donation %s: S3 client not configured", donation.ID.String())
		return
	}

	// Get campaign title
	campaignTitle, err := s.campaignService.GetCampaignTitle(ctx, donation.CampaignID)
	if err != nil {
		log.Printf("Error getting campaign title for receipt generation: %v", err)
		campaignTitle = "Campaña sin título"
	}

	// Get donor information
	donorInfo, err := s.donorService.GetDonor(ctx, donation.DonorID)
	var donorName string
	if err != nil {
		log.Printf("Error getting donor info for receipt generation: %v", err)
		donorName = "Donante"
	} else {
		donorName = fmt.Sprintf("%s %s", donorInfo.FirstName, donorInfo.LastName)
	}

	// Get payment method name
	paymentMethodName := "Método de pago"
	if donation.PaymentMethod != nil {
		paymentMethodName = donation.PaymentMethod.Name
	}

	// Prepare receipt data
	receiptData := receipt.ReceiptData{
		DonationID:    donation.ID,
		CampaignTitle: campaignTitle,
		DonorName:     donorName,
		Amount:        donation.Amount,
		Date:          donation.Date,
		PaymentMethod: paymentMethodName,
		IsAnonymous:   donation.IsAnonymous,
	}

	// Generate PDF
	pdfBytes, err := s.pdfGenerator.Generate(receiptData)
	if err != nil {
		log.Printf("Error generating PDF receipt for donation %s: %v", donation.ID.String(), err)
		return
	}

	// Upload to S3
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("receipt-%d.pdf", timestamp)
	
	uploadResp, err := s.s3Client.UploadBytes(ctx, s3client.UploadBytesRequest{
		Data:         pdfBytes,
		FileName:     fileName,
		ResourceType: "donation/receipts",
		ResourceID:   donation.ID.String(),
	})
	if err != nil {
		log.Printf("Error uploading receipt to S3 for donation %s: %v", donation.ID.String(), err)
		return
	}

	// Update donation with receipt URL
	if err := s.repo.UpdateReceiptURL(ctx, donation.ID, uploadResp.URL); err != nil {
		log.Printf("Error updating receipt URL for donation %s: %v", donation.ID.String(), err)
		return
	}

	log.Printf("✅ Receipt generated and uploaded successfully for donation %s", donation.ID.String())
}
