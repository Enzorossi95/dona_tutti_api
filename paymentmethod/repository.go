package paymentmethod

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetPaymentMethods(ctx context.Context) ([]PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, id int) (PaymentMethod, error)
	GetPaymentMethodByCode(ctx context.Context, code string) (PaymentMethod, error)

	GetCampaignPaymentMethods(ctx context.Context, campaignID uuid.UUID) ([]CampaignPaymentMethod, error)
	GetCampaignPaymentMethod(ctx context.Context, id int) (CampaignPaymentMethod, error)
	CreateCampaignPaymentMethod(ctx context.Context, cpm CampaignPaymentMethod) (int, error)
	UpdateCampaignPaymentMethod(ctx context.Context, cpm CampaignPaymentMethod) error
	DeleteCampaignPaymentMethod(ctx context.Context, id int) error

	CreateTransferDetail(ctx context.Context, detail TransferDetail) error
	UpdateTransferDetail(ctx context.Context, detail TransferDetail) error
	DeleteTransferDetail(ctx context.Context, campaignPaymentMethodID int) error

	CreateCashLocations(ctx context.Context, locations []CashLocation) error
	UpdateCashLocation(ctx context.Context, location CashLocation) error
	DeleteCashLocationsByCampaignPaymentMethod(ctx context.Context, campaignPaymentMethodID int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Payment Methods
func (r *repository) GetPaymentMethods(ctx context.Context) ([]PaymentMethod, error) {
	var models []PaymentMethodModel

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name").
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	entities := make([]PaymentMethod, len(models))
	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities, nil
}

func (r *repository) GetPaymentMethod(ctx context.Context, id int) (PaymentMethod, error) {
	var model PaymentMethodModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return PaymentMethod{}, fmt.Errorf("payment method with id %d not found", id)
		}
		return PaymentMethod{}, fmt.Errorf("failed to get payment method: %w", err)
	}

	return model.ToEntity(), nil
}

func (r *repository) GetPaymentMethodByCode(ctx context.Context, code string) (PaymentMethod, error) {
	var model PaymentMethodModel

	err := r.db.WithContext(ctx).Where("code = ?", code).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return PaymentMethod{}, fmt.Errorf("payment method with code %s not found", code)
		}
		return PaymentMethod{}, fmt.Errorf("failed to get payment method: %w", err)
	}

	return model.ToEntity(), nil
}

// Campaign Payment Methods
func (r *repository) GetCampaignPaymentMethods(ctx context.Context, campaignID uuid.UUID) ([]CampaignPaymentMethod, error) {
	var models []CampaignPaymentMethodModel

	err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Preload("TransferDetails").
		Preload("CashLocations").
		Where("campaign_id = ? AND is_active = ?", campaignID, true).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign payment methods: %w", err)
	}

	entities := make([]CampaignPaymentMethod, len(models))
	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities, nil
}

func (r *repository) GetCampaignPaymentMethod(ctx context.Context, id int) (CampaignPaymentMethod, error) {
	var model CampaignPaymentMethodModel

	err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Preload("TransferDetails").
		Preload("CashLocations").
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return CampaignPaymentMethod{}, fmt.Errorf("campaign payment method with id %d not found", id)
		}
		return CampaignPaymentMethod{}, fmt.Errorf("failed to get campaign payment method: %w", err)
	}

	return model.ToEntity(), nil
}

func (r *repository) CreateCampaignPaymentMethod(ctx context.Context, cpm CampaignPaymentMethod) (int, error) {
	var model CampaignPaymentMethodModel
	model.FromEntity(cpm)

	result := r.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create campaign payment method: %w", result.Error)
	}

	return model.ID, nil
}

func (r *repository) UpdateCampaignPaymentMethod(ctx context.Context, cpm CampaignPaymentMethod) error {
	var model CampaignPaymentMethodModel
	model.FromEntity(cpm)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return fmt.Errorf("failed to update campaign payment method: %w", err)
	}

	return nil
}

func (r *repository) DeleteCampaignPaymentMethod(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&CampaignPaymentMethodModel{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete campaign payment method: %w", err)
	}

	return nil
}

// Transfer Details
func (r *repository) CreateTransferDetail(ctx context.Context, detail TransferDetail) error {
	var model TransferDetailModel
	model.FromEntity(detail)

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return fmt.Errorf("failed to create transfer detail: %w", err)
	}

	return nil
}

func (r *repository) UpdateTransferDetail(ctx context.Context, detail TransferDetail) error {
	var model TransferDetailModel
	model.FromEntity(detail)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return fmt.Errorf("failed to update transfer detail: %w", err)
	}

	return nil
}

func (r *repository) DeleteTransferDetail(ctx context.Context, campaignPaymentMethodID int) error {
	err := r.db.WithContext(ctx).
		Where("campaign_payment_method_id = ?", campaignPaymentMethodID).
		Delete(&TransferDetailModel{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete transfer detail: %w", err)
	}

	return nil
}

// Cash Locations
func (r *repository) CreateCashLocations(ctx context.Context, locations []CashLocation) error {
	if len(locations) == 0 {
		return nil
	}

	models := make([]CashLocationModel, len(locations))
	for i, location := range locations {
		models[i].FromEntity(location)
	}

	err := r.db.WithContext(ctx).Create(&models).Error
	if err != nil {
		return fmt.Errorf("failed to create cash locations: %w", err)
	}

	return nil
}

func (r *repository) UpdateCashLocation(ctx context.Context, location CashLocation) error {
	var model CashLocationModel
	model.FromEntity(location)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return fmt.Errorf("failed to update cash location: %w", err)
	}

	return nil
}

func (r *repository) DeleteCashLocationsByCampaignPaymentMethod(ctx context.Context, campaignPaymentMethodID int) error {
	err := r.db.WithContext(ctx).
		Where("campaign_payment_method_id = ?", campaignPaymentMethodID).
		Delete(&CashLocationModel{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete cash locations: %w", err)
	}

	return nil
}
