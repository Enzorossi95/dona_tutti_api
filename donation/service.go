package donation

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetDonation(ctx context.Context, id uuid.UUID) (Donation, error)
	CreateDonation(ctx context.Context, donation Donation) (uuid.UUID, error)
	UpdateDonation(ctx context.Context, donation Donation) error
	ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error)
}

type service struct {
	repo DonationRepository
}

func NewService(repo DonationRepository) Service {
	return &service{repo: repo}
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
	return donation.ID, nil
}

func (s *service) UpdateDonation(ctx context.Context, donation Donation) error {
	return s.repo.UpdateDonation(ctx, donation)
}

func (s *service) ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error) {
	return s.repo.ListDonationsByCampaign(ctx, campaignID)
}
