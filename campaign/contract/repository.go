package contract

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for campaign contract data access
type Repository interface {
	Create(ctx context.Context, contract CampaignContract) error
	Update(ctx context.Context, contract CampaignContract) error
	GetByCampaignID(ctx context.Context, campaignID uuid.UUID) (CampaignContract, error)
	GetByID(ctx context.Context, id uuid.UUID) (CampaignContract, error)
	ExistsByCampaignID(ctx context.Context, campaignID uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of the contract repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create saves a new campaign contract to the database
func (r *repository) Create(ctx context.Context, contract CampaignContract) error {
	var model CampaignContractModel
	model.FromEntity(contract)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create campaign contract: %w", err)
	}

	return nil
}

// Update updates an existing campaign contract in the database
func (r *repository) Update(ctx context.Context, contract CampaignContract) error {
	var model CampaignContractModel
	model.FromEntity(contract)

	if err := r.db.WithContext(ctx).
		Model(&CampaignContractModel{}).
		Where("campaign_id = ?", contract.CampaignID).
		Updates(map[string]interface{}{
			"accepted_at":           model.AcceptedAt,
			"acceptance_ip":         model.AcceptanceIP,
			"acceptance_user_agent": model.AcceptanceUserAgent,
		}).Error; err != nil {
		return fmt.Errorf("failed to update campaign contract: %w", err)
	}

	return nil
}

// GetByCampaignID retrieves a contract by campaign ID
func (r *repository) GetByCampaignID(ctx context.Context, campaignID uuid.UUID) (CampaignContract, error) {
	var model CampaignContractModel

	err := r.db.WithContext(ctx).
		Where("campaign_id = ?", campaignID).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return CampaignContract{}, fmt.Errorf("contract not found for campaign %s", campaignID)
		}
		return CampaignContract{}, fmt.Errorf("failed to get contract: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByID retrieves a contract by its ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (CampaignContract, error) {
	var model CampaignContractModel

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return CampaignContract{}, fmt.Errorf("contract not found with id %s", id)
		}
		return CampaignContract{}, fmt.Errorf("failed to get contract: %w", err)
	}

	return model.ToEntity(), nil
}

// ExistsByCampaignID checks if a contract exists for a campaign
func (r *repository) ExistsByCampaignID(ctx context.Context, campaignID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&CampaignContractModel{}).
		Where("campaign_id = ?", campaignID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check contract existence: %w", err)
	}

	return count > 0, nil
}

