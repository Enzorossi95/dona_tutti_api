package campaign

import (
	"context"
	"fmt"
	"time"

	apierrors "dona_tutti_api/errors"
	"dona_tutti_api/organizer"
	"dona_tutti_api/paymentmethod"

	"github.com/google/uuid"
)

// CampaignInfo represents minimal campaign information for external packages
type CampaignInfo struct {
	ID          uuid.UUID
	Title       string
	Goal        float64
	OrganizerID uuid.UUID
	Status      string
}

type Service interface {
	GetCampaign(ctx context.Context, id uuid.UUID) (Campaign, error)
	ListCampaigns(ctx context.Context) ([]Campaign, error)
	CreateCampaign(ctx context.Context, campaign Campaign) (uuid.UUID, error)
	UpdateCampaignImage(ctx context.Context, id uuid.UUID, imageURL string) error
	UpdateStatus(ctx context.Context, campaignID uuid.UUID, status string) error
	GetCampaignTitle(ctx context.Context, campaignID uuid.UUID) (string, error)
	GetCampaignInfo(ctx context.Context, campaignID uuid.UUID) (CampaignInfo, error)
	GetSummary(ctx context.Context) (Summary, error)
}

type CampaignRepository interface {
	GetCampaign(ctx context.Context, id uuid.UUID) (Campaign, error)
	ListCampaigns(ctx context.Context) ([]Campaign, error)
	CreateCampaign(ctx context.Context, campaign Campaign) error
	UpdateCampaignImage(ctx context.Context, id uuid.UUID, imageURL string) error
	UpdateStatus(ctx context.Context, campaignID uuid.UUID, status string) error
	GetSummary(ctx context.Context) (Summary, error)
}

type PaymentMethodService interface {
	CreateCampaignPaymentMethod(ctx context.Context, req paymentmethod.CreateCampaignPaymentMethodRequest) (int, error)
}

type OrganizerService interface {
	GetOrganizer(ctx context.Context, id uuid.UUID) (organizer.Organizer, error)
	UpdateOrganizer(ctx context.Context, organizer organizer.Organizer) error
}

type service struct {
	repo             CampaignRepository
	paymentMethodSvc PaymentMethodService
	organizerSvc     OrganizerService
}

func NewService(repo CampaignRepository, paymentMethodSvc PaymentMethodService, organizerSvc OrganizerService) Service {
	return &service{repo: repo, paymentMethodSvc: paymentMethodSvc, organizerSvc: organizerSvc}
}

func (s *service) GetCampaign(ctx context.Context, id uuid.UUID) (Campaign, error) {
	return s.repo.GetCampaign(ctx, id)
}

func (s *service) ListCampaigns(ctx context.Context) ([]Campaign, error) {
	return s.repo.ListCampaigns(ctx)
}

func (s *service) CreateCampaign(ctx context.Context, campaign Campaign) (uuid.UUID, error) {
	// Generate new ID and set timestamps
	campaign.ID = uuid.New()
	campaign.CreatedAt = time.Now()

	// Set default status to draft for new campaigns
	if campaign.Status == "" {
		campaign.Status = StatusDraft
	}

	// Validate status
	if !IsValidStatus(campaign.Status) {
		return uuid.Nil, apierrors.NewFieldValidationError("status", "invalid campaign status")
	}

	// Validate required fields
	if campaign.Title == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("title", "campaign title is required")
	}
	if campaign.Description == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("description", "campaign description is required")
	}
	if campaign.Goal <= 0 {
		return uuid.Nil, apierrors.NewFieldValidationError("goal", "campaign goal must be greater than 0")
	}
	if campaign.Urgency < 1 || campaign.Urgency > 10 {
		return uuid.Nil, apierrors.NewFieldValidationError("urgency", "campaign urgency must be between 1 and 10")
	}
	if campaign.StartDate.IsZero() {
		campaign.StartDate = time.Now()
	}
	if campaign.EndDate.IsZero() || campaign.EndDate.Before(campaign.StartDate) {
		return uuid.Nil, apierrors.NewFieldValidationError("end_date", "campaign end date must be after start date")
	}

	// Validate and update organizer information
	if campaign.Organizer == nil {
		return uuid.Nil, apierrors.NewFieldValidationError("organizer", "organizer information is required")
	}
	if campaign.Organizer.ID == uuid.Nil {
		return uuid.Nil, apierrors.NewFieldValidationError("organizer.id", "organizer ID is required")
	}
	if campaign.Organizer.Name == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("organizer.name", "organizer name is required")
	}
	if campaign.Organizer.Email == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("organizer.email", "organizer email is required")
	}
	if campaign.Organizer.Phone == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("organizer.phone", "organizer phone is required")
	}

	// Get existing organizer to preserve UserID and other fields
	existingOrganizer, err := s.organizerSvc.GetOrganizer(ctx, campaign.Organizer.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get existing organizer: %w", err)
	}

	// Update only the provided fields, preserving existing data
	updatedOrganizer := existingOrganizer
	updatedOrganizer.Name = campaign.Organizer.Name
	updatedOrganizer.Email = campaign.Organizer.Email
	updatedOrganizer.Phone = campaign.Organizer.Phone

	// Update organizer information
	if err := s.organizerSvc.UpdateOrganizer(ctx, updatedOrganizer); err != nil {
		return uuid.Nil, fmt.Errorf("failed to update organizer: %w", err)
	}

	if err := s.repo.CreateCampaign(ctx, campaign); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	if len(campaign.PaymentMethods) > 0 {
		for _, paymentMethod := range campaign.PaymentMethods {
			req := paymentmethod.CreateCampaignPaymentMethodRequest{
				CampaignID:      campaign.ID,
				PaymentMethodID: paymentMethod.PaymentMethodID,
				Instructions:    paymentMethod.Instructions,
			}

			if _, err := s.paymentMethodSvc.CreateCampaignPaymentMethod(ctx, req); err != nil {
				return uuid.Nil, fmt.Errorf("failed to create campaign payment method: %w", err)
			}
		}
	}

	return campaign.ID, nil
}

func (s *service) UpdateCampaignImage(ctx context.Context, id uuid.UUID, imageURL string) error {
	// Check if campaign exists
	_, err := s.repo.GetCampaign(ctx, id)
	if err != nil {
		return fmt.Errorf("campaign not found: %w", err)
	}

	// Update campaign image
	if err := s.repo.UpdateCampaignImage(ctx, id, imageURL); err != nil {
		return fmt.Errorf("failed to update campaign image: %w", err)
	}

	return nil
}

func (s *service) GetSummary(ctx context.Context) (Summary, error) {
	return s.repo.GetSummary(ctx)
}

func (s *service) UpdateStatus(ctx context.Context, campaignID uuid.UUID, status string) error {
	// Validate new status
	if !IsValidStatus(status) {
		return apierrors.NewFieldValidationError("status", "invalid campaign status")
	}

	// Get current campaign
	campaign, err := s.repo.GetCampaign(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("campaign not found: %w", err)
	}

	// Validate status transition
	if !CanTransitionTo(campaign.Status, status) {
		return apierrors.NewFieldValidationError("status",
			fmt.Sprintf("cannot transition from %s to %s", campaign.Status, status))
	}

	// Additional validation: cannot transition to active without contract
	if status == StatusActive && campaign.Status == StatusPendingApproval {
		// This validation will be handled by the admin approval flow
		// For now, we just allow the transition
	}

	// Update status
	if err := s.repo.UpdateStatus(ctx, campaignID, status); err != nil {
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	return nil
}

func (s *service) GetCampaignTitle(ctx context.Context, campaignID uuid.UUID) (string, error) {
	campaign, err := s.repo.GetCampaign(ctx, campaignID)
	if err != nil {
		return "", fmt.Errorf("campaign not found: %w", err)
	}
	return campaign.Title, nil
}

func (s *service) GetCampaignInfo(ctx context.Context, campaignID uuid.UUID) (CampaignInfo, error) {
	campaign, err := s.repo.GetCampaign(ctx, campaignID)
	if err != nil {
		return CampaignInfo{}, fmt.Errorf("campaign not found: %w", err)
	}

	return CampaignInfo{
		ID:          campaign.ID,
		Title:       campaign.Title,
		Goal:        campaign.Goal,
		OrganizerID: campaign.OrganizerID,
		Status:      campaign.Status,
	}, nil
}
