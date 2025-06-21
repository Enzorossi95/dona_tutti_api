package donation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DonationRepository interface {
	GetDonation(ctx context.Context, id uuid.UUID) (Donation, error)
	CreateDonation(ctx context.Context, donation Donation) error
	UpdateDonation(ctx context.Context, donation Donation) error
	ListDonations(ctx context.Context) ([]Donation, error)
	ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error)
}

type donationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) DonationRepository {
	return &donationRepository{db: db}
}

func (r *donationRepository) GetDonation(ctx context.Context, id uuid.UUID) (Donation, error) {
	var model DonationModel
	if err := r.db.WithContext(ctx).Preload("Donor").Where("id = ?", id).First(&model).Error; err != nil {
		return Donation{}, fmt.Errorf("failed to get donation: %w", err)
	}
	return model.ToEntity(), nil
}

func (r *donationRepository) CreateDonation(ctx context.Context, donation Donation) error {
	model := DonationModel{}
	model.FromEntity(donation)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create donation: %w", err)
	}
	return nil
}

func (r *donationRepository) UpdateDonation(ctx context.Context, donation Donation) error {
	model := DonationModel{}
	model.FromEntity(donation)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("failed to update donation: %w", err)
	}
	return nil
}

func (r *donationRepository) ListDonations(ctx context.Context) ([]Donation, error) {
	var models []DonationModel
	if err := r.db.WithContext(ctx).Preload("Donor").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list donations: %w", err)
	}
	donations := make([]Donation, len(models))
	for i, model := range models {
		donations[i] = model.ToEntity()
	}
	return donations, nil
}

func (r *donationRepository) ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error) {
	var models []DonationModel
	if err := r.db.WithContext(ctx).Preload("Donor").Where("campaign_id = ?", campaignID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list donations by campaign: %w", err)
	}
	donations := make([]Donation, len(models))
	for i, model := range models {
		donations[i] = model.ToEntity()
	}
	return donations, nil
}
