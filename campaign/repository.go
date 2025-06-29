package campaign

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type campaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *campaignRepository {
	return &campaignRepository{
		db: db,
	}
}

func (r *campaignRepository) GetCampaign(ctx context.Context, id uuid.UUID) (Campaign, error) {
	var campaignModel CampaignModel

	// Get the campaign model
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&campaignModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Campaign{}, fmt.Errorf("campaign with id %s not found", id.String())
		}
		return Campaign{}, fmt.Errorf("failed to get campaign: %w", err)
	}

	var paymentMethods []CampaignPaymentMethodModel
	err = r.db.WithContext(ctx).
		Table("campaign_payment_methods cpm").
		Select("cpm.id, cpm.payment_method_id, pm.code, pm.name, cpm.instructions").
		Joins("JOIN payment_methods pm ON cpm.payment_method_id = pm.id").
		Where("cpm.campaign_id = ? AND cpm.is_active = ? AND pm.is_active = ?", id, true, true).
		Scan(&paymentMethods).Error
	if err != nil {
		fmt.Printf("Warning: failed to get payment methods for campaign %s: %v\n", id.String(), err)
	}

	campaignModel.PaymentMethods = paymentMethods
	return campaignModel.ToEntity(), nil
}

func (r *campaignRepository) ListCampaigns(ctx context.Context) ([]Campaign, error) {
	type Campaigns struct {
		CampaignModel
	}

	var campaignModels []Campaigns

	err := r.db.WithContext(ctx).
		Table("campaigns").
		Select(`campaigns.*`).
		Order("campaigns.created_at DESC").
		Scan(&campaignModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns: %w", err)
	}

	// Convert to domain entities
	campaigns := make([]Campaign, len(campaignModels))
	for i, model := range campaignModels {
		campaigns[i] = model.CampaignModel.ToEntity()
	}

	return campaigns, nil
}

func (r *campaignRepository) CreateCampaign(ctx context.Context, campaign Campaign) error {

	// Convert domain entity to database model
	var campaignModel CampaignModel
	campaignModel.FromEntity(campaign)

	// Create the campaign
	err := r.db.WithContext(ctx).Create(&campaignModel).Error
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	return nil
}

func (r *campaignRepository) GetSummary(ctx context.Context) (Summary, error) {
	var summary Summary

	err := r.db.WithContext(ctx).
		Raw(`
                       SELECT
                               COUNT(DISTINCT c.id) AS total_campaigns,
                               COALESCE(SUM(c.goal), 0) AS total_goal,
                               COUNT(DISTINCT d.donor_id) AS total_contributors
                       FROM campaigns c
                       LEFT JOIN donations d ON c.id = d.campaign_id
               `).
		Scan(&summary).Error

	if err != nil {
		return Summary{}, fmt.Errorf(" to get summary: %w", err)
	}

	return summary, nil
}
