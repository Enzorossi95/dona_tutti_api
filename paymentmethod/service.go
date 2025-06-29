package paymentmethod

import (
	"context"
	"fmt"

	apierrors "dona_tutti_api/errors"

	"github.com/google/uuid"
)

type Service interface {
	GetPaymentMethods(ctx context.Context) ([]PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, id int) (PaymentMethod, error)

	GetCampaignPaymentMethods(ctx context.Context, campaignID uuid.UUID) ([]CampaignPaymentMethod, error)
	CreateCampaignPaymentMethod(ctx context.Context, req CreateCampaignPaymentMethodRequest) (int, error)
	UpdateCampaignPaymentMethod(ctx context.Context, id int, req CreateCampaignPaymentMethodRequest) error
	DeleteCampaignPaymentMethod(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPaymentMethods(ctx context.Context) ([]PaymentMethod, error) {
	return s.repo.GetPaymentMethods(ctx)
}

func (s *service) GetPaymentMethod(ctx context.Context, id int) (PaymentMethod, error) {
	return s.repo.GetPaymentMethod(ctx, id)
}

func (s *service) GetCampaignPaymentMethods(ctx context.Context, campaignID uuid.UUID) ([]CampaignPaymentMethod, error) {
	return s.repo.GetCampaignPaymentMethods(ctx, campaignID)
}

func (s *service) CreateCampaignPaymentMethod(ctx context.Context, req CreateCampaignPaymentMethodRequest) (int, error) {
	// Validate that the payment method exists
	paymentMethod, err := s.repo.GetPaymentMethod(ctx, req.PaymentMethodID)
	if err != nil {
		return 0, apierrors.NewFieldValidationError("payment_method_id", "invalid payment method")
	}

	// Create the campaign payment method
	cpm := CampaignPaymentMethod{
		CampaignID:      req.CampaignID,
		PaymentMethodID: req.PaymentMethodID,
		Instructions:    req.Instructions,
		IsActive:        true,
	}

	createdID, err := s.repo.CreateCampaignPaymentMethod(ctx, cpm)
	if err != nil {
		return 0, fmt.Errorf("failed to create campaign payment method: %w", err)
	}

	if createdID == 0 {
		return 0, fmt.Errorf("failed to find created campaign payment method")
	}

	// Handle transfer details if payment method is transfer
	if paymentMethod.Code == "transfer" && req.TransferDetails != nil {
		transferDetail := TransferDetail{
			CampaignPaymentMethodID: createdID,
			BankName:                req.TransferDetails.BankName,
			AccountHolder:           req.TransferDetails.AccountHolder,
			CBU:                     req.TransferDetails.CBU,
			Alias:                   req.TransferDetails.Alias,
			SwiftCode:               req.TransferDetails.SwiftCode,
			AdditionalNotes:         req.TransferDetails.AdditionalNotes,
		}

		err = s.repo.CreateTransferDetail(ctx, transferDetail)
		if err != nil {
			// Rollback: delete the campaign payment method
			s.repo.DeleteCampaignPaymentMethod(ctx, createdID)
			return 0, fmt.Errorf("failed to create transfer details: %w", err)
		}
	}

	// Handle cash locations if payment method is cash
	if paymentMethod.Code == "cash" && len(req.CashLocations) > 0 {
		cashLocations := make([]CashLocation, len(req.CashLocations))
		for i, location := range req.CashLocations {
			cashLocations[i] = CashLocation{
				CampaignPaymentMethodID: createdID,
				LocationName:            location.LocationName,
				Address:                 location.Address,
				ContactInfo:             location.ContactInfo,
				AvailableHours:          location.AvailableHours,
				AdditionalNotes:         location.AdditionalNotes,
			}
		}

		err = s.repo.CreateCashLocations(ctx, cashLocations)
		if err != nil {
			// Rollback: delete the campaign payment method and transfer details
			s.repo.DeleteTransferDetail(ctx, createdID)
			s.repo.DeleteCampaignPaymentMethod(ctx, createdID)
			return 0, fmt.Errorf("failed to create cash locations: %w", err)
		}
	}

	return createdID, nil
}

func (s *service) UpdateCampaignPaymentMethod(ctx context.Context, id int, req CreateCampaignPaymentMethodRequest) error {
	// Get existing campaign payment method
	existing, err := s.repo.GetCampaignPaymentMethod(ctx, id)
	if err != nil {
		return err
	}

	// Get payment method info
	paymentMethod, err := s.repo.GetPaymentMethod(ctx, req.PaymentMethodID)
	if err != nil {
		return apierrors.NewFieldValidationError("payment_method_id", "invalid payment method")
	}

	// Update the campaign payment method
	existing.Instructions = req.Instructions
	err = s.repo.UpdateCampaignPaymentMethod(ctx, existing)
	if err != nil {
		return fmt.Errorf("failed to update campaign payment method: %w", err)
	}

	// Handle transfer details
	if paymentMethod.Code == "transfer" {
		//----- TODO: No es necesario eliminar los transfer details, se debe actualizar el transfer detail existente
		// Delete existing transfer details
		s.repo.DeleteTransferDetail(ctx, id)

		// Create new ones if provided
		if req.TransferDetails != nil {
			transferDetail := TransferDetail{
				CampaignPaymentMethodID: id,
				BankName:                req.TransferDetails.BankName,
				AccountHolder:           req.TransferDetails.AccountHolder,
				CBU:                     req.TransferDetails.CBU,
				Alias:                   req.TransferDetails.Alias,
				SwiftCode:               req.TransferDetails.SwiftCode,
				AdditionalNotes:         req.TransferDetails.AdditionalNotes,
			}

			err = s.repo.CreateTransferDetail(ctx, transferDetail)
			if err != nil {
				return fmt.Errorf("failed to update transfer details: %w", err)
			}
		}
	}

	// Handle cash locations
	if paymentMethod.Code == "cash" {
		//----- TODO: No es necesario eliminar las cash locations, se debe actualizar la cash location existente
		// Delete existing cash locations
		s.repo.DeleteCashLocationsByCampaignPaymentMethod(ctx, id)

		// Create new ones if provided
		if len(req.CashLocations) > 0 {
			cashLocations := make([]CashLocation, len(req.CashLocations))
			for i, location := range req.CashLocations {
				cashLocations[i] = CashLocation{
					CampaignPaymentMethodID: id,
					LocationName:            location.LocationName,
					Address:                 location.Address,
					ContactInfo:             location.ContactInfo,
					AvailableHours:          location.AvailableHours,
					AdditionalNotes:         location.AdditionalNotes,
				}
			}

			err = s.repo.CreateCashLocations(ctx, cashLocations)
			if err != nil {
				return fmt.Errorf("failed to update cash locations: %w", err)
			}
		}
	}

	return nil
}

func (s *service) DeleteCampaignPaymentMethod(ctx context.Context, id int) error {
	// Delete associated details first (handled by foreign key constraints with CASCADE)
	return s.repo.DeleteCampaignPaymentMethod(ctx, id)
}
