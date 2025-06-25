package campaign

import (
	"context"
	"fmt"
	"time"

	apierrors "dona_tutti_api/errors"

	"github.com/google/uuid"
)

type Service interface {
	GetCampaign(ctx context.Context, id uuid.UUID) (Campaign, error)
	ListCampaigns(ctx context.Context) ([]Campaign, error)
	CreateCampaign(ctx context.Context, campaign Campaign) (uuid.UUID, error)
	GetSummary(ctx context.Context) (Summary, error)
}

type CampaignRepository interface {
	GetCampaign(ctx context.Context, id uuid.UUID) (Campaign, error)
	ListCampaigns(ctx context.Context) ([]Campaign, error)
	CreateCampaign(ctx context.Context, campaign Campaign) error
	GetSummary(ctx context.Context) (Summary, error)
}

type service struct {
	repo CampaignRepository
}

func NewService(repo CampaignRepository) Service {
	return &service{repo: repo}
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

	// Set default status if not provided
	if campaign.Status == "" {
		campaign.Status = "active"
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

	if err := s.repo.CreateCampaign(ctx, campaign); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign.ID, nil
}

func (s *service) GetSummary(ctx context.Context) (Summary, error) {
	return s.repo.GetSummary(ctx)
}
