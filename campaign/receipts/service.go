package receipts

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetReceiptsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Receipt, error)
	GetReceipt(ctx context.Context, id uuid.UUID) (Receipt, error)
	CreateReceipt(ctx context.Context, campaignID uuid.UUID, req ReceiptCreateRequest) (Receipt, error)
	UpdateReceipt(ctx context.Context, id uuid.UUID, req ReceiptUpdateRequest) (Receipt, error)
	DeleteReceipt(ctx context.Context, id uuid.UUID) error
	UpdateDocumentURL(ctx context.Context, id uuid.UUID, documentURL string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetReceiptsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Receipt, error) {
	return s.repo.GetReceiptsByCampaign(ctx, campaignID)
}

func (s *service) GetReceipt(ctx context.Context, id uuid.UUID) (Receipt, error) {
	receipt, err := s.repo.GetReceipt(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Receipt{}, errors.New("receipt not found")
		}
		return Receipt{}, err
	}
	return receipt, nil
}

func (s *service) CreateReceipt(ctx context.Context, campaignID uuid.UUID, req ReceiptCreateRequest) (Receipt, error) {
	// Set default quantity if not provided
	quantity := req.Quantity
	if quantity == 0 {
		quantity = 1
	}

	receipt := Receipt{
		ID:          uuid.New(),
		CampaignID:  campaignID,
		Provider:    req.Provider,
		Name:        req.Name,
		Description: req.Description,
		Total:       req.Total,
		Quantity:    quantity,
		Date:        req.Date,
		Note:        req.Note,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateReceipt(ctx, receipt); err != nil {
		return Receipt{}, err
	}

	return receipt, nil
}

func (s *service) UpdateReceipt(ctx context.Context, id uuid.UUID, req ReceiptUpdateRequest) (Receipt, error) {
	// Get existing receipt
	receipt, err := s.repo.GetReceipt(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Receipt{}, errors.New("receipt not found")
		}
		return Receipt{}, err
	}

	// Update fields if provided
	if req.Provider != nil {
		receipt.Provider = *req.Provider
	}
	if req.Name != nil {
		receipt.Name = *req.Name
	}
	if req.Description != nil {
		receipt.Description = *req.Description
	}
	if req.Total != nil {
		receipt.Total = *req.Total
	}
	if req.Quantity != nil {
		receipt.Quantity = *req.Quantity
	}
	if req.Date != nil {
		receipt.Date = *req.Date
	}
	if req.Note != nil {
		receipt.Note = req.Note
	}

	receipt.UpdatedAt = time.Now()

	if err := s.repo.UpdateReceipt(ctx, receipt); err != nil {
		return Receipt{}, err
	}

	return receipt, nil
}

func (s *service) DeleteReceipt(ctx context.Context, id uuid.UUID) error {
	// Check if receipt exists
	_, err := s.repo.GetReceipt(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("receipt not found")
		}
		return err
	}

	return s.repo.DeleteReceipt(ctx, id)
}

func (s *service) UpdateDocumentURL(ctx context.Context, id uuid.UUID, documentURL string) error {
	// Check if receipt exists
	_, err := s.repo.GetReceipt(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("receipt not found")
		}
		return err
	}

	return s.repo.UpdateDocumentURL(ctx, id, documentURL)
}