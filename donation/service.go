package donation

import (
	"context"
	"dona_tutti_api/donor"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetDonation(ctx context.Context, id uuid.UUID) (Donation, error)
	CreateDonation(ctx context.Context, donation Donation) (uuid.UUID, error)
	CreateDonationWithRequest(ctx context.Context, campaignID uuid.UUID, req CreateDonationRequest) (uuid.UUID, error)
	UpdateDonation(ctx context.Context, donation Donation) error
	ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error)
}

type service struct {
	repo        DonationRepository
	donorService donor.Service
}

func NewService(repo DonationRepository, donorService donor.Service) Service {
	return &service{
		repo:         repo,
		donorService: donorService,
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
	return donation.ID, nil
}

func (s *service) UpdateDonation(ctx context.Context, donation Donation) error {
	return s.repo.UpdateDonation(ctx, donation)
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
