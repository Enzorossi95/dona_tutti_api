package receipts

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetReceiptsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Receipt, error)
	GetReceipt(ctx context.Context, id uuid.UUID) (Receipt, error)
	CreateReceipt(ctx context.Context, receipt Receipt) error
	UpdateReceipt(ctx context.Context, receipt Receipt) error
	DeleteReceipt(ctx context.Context, id uuid.UUID) error
	UpdateDocumentURL(ctx context.Context, id uuid.UUID, documentURL string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetReceiptsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Receipt, error) {
	var models []ReceiptModel
	if err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Order("date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	receipts := make([]Receipt, len(models))
	for i, model := range models {
		receipts[i] = model.ToEntity()
	}

	return receipts, nil
}

func (r *repository) GetReceipt(ctx context.Context, id uuid.UUID) (Receipt, error) {
	var model ReceiptModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return Receipt{}, err
	}

	return model.ToEntity(), nil
}

func (r *repository) CreateReceipt(ctx context.Context, receipt Receipt) error {
	var model ReceiptModel
	model.FromEntity(receipt)

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *repository) UpdateReceipt(ctx context.Context, receipt Receipt) error {
	var model ReceiptModel
	model.FromEntity(receipt)

	return r.db.WithContext(ctx).Where("id = ?", receipt.ID).Updates(&model).Error
}

func (r *repository) DeleteReceipt(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ReceiptModel{}).Error
}

func (r *repository) UpdateDocumentURL(ctx context.Context, id uuid.UUID, documentURL string) error {
	return r.db.WithContext(ctx).Model(&ReceiptModel{}).Where("id = ?", id).Update("document_url", documentURL).Error
}