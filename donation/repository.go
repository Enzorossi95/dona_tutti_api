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
	if err := r.db.WithContext(ctx).
		Preload("Donor").
		Where("id = ?", id).
		First(&model).Error; err != nil {
		return Donation{}, fmt.Errorf("failed to get donation: %w", err)
	}

	// Get payment method info
	var paymentMethod PaymentMethodModel
	if err := r.db.WithContext(ctx).
		Table("payment_methods").
		Select("id, code, name").
		Where("id = ?", model.PaymentMethodID).
		First(&paymentMethod).Error; err == nil {
		model.PaymentMethod = &paymentMethod
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

func (r *donationRepository) ListDonationsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Donation, error) {
	var models []DonationModel
	if err := r.db.WithContext(ctx).Preload("Donor").Where("campaign_id = ?", campaignID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list donations by campaign: %w", err)
	}

	// Get payment method info for all donations
	r.loadPaymentMethods(ctx, models)

	donations := make([]Donation, len(models))
	for i, model := range models {
		donations[i] = model.ToEntity()
	}
	return donations, nil
}

// loadPaymentMethods loads payment method information for a slice of donation models
func (r *donationRepository) loadPaymentMethods(ctx context.Context, models []DonationModel) {
	if len(models) == 0 {
		return
	}

	// Get unique payment method IDs
	paymentMethodIDs := make([]int, 0)
	paymentMethodMap := make(map[int]bool)
	for _, model := range models {
		if !paymentMethodMap[model.PaymentMethodID] {
			paymentMethodIDs = append(paymentMethodIDs, model.PaymentMethodID)
			paymentMethodMap[model.PaymentMethodID] = true
		}
	}

	// Load payment methods in batch
	var paymentMethods []PaymentMethodModel
	if err := r.db.WithContext(ctx).
		Table("payment_methods").
		Select("id, code, name").
		Where("id IN ?", paymentMethodIDs).
		Find(&paymentMethods).Error; err != nil {
		return // Don't fail if we can't load payment methods
	}

	// Create map for quick lookup
	pmMap := make(map[int]*PaymentMethodModel)
	for i := range paymentMethods {
		pmMap[paymentMethods[i].ID] = &paymentMethods[i]
	}

	// Assign payment methods to models
	for i := range models {
		if pm, exists := pmMap[models[i].PaymentMethodID]; exists {
			models[i].PaymentMethod = pm
		}
	}
}
