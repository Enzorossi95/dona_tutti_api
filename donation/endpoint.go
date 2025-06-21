package donation

import (
	"context"
	"errors"
	"fmt"
	"time"

	apierrors "microservice_go/errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type GetDonationRequestModel struct {
	ID uuid.UUID `json:"id"`
}

type GetDonationResponseModel struct {
	Donation Donation `json:"donation"`
}

type ListDonationsRequestModel struct {
	CampaignID *uuid.UUID `json:"campaign_id,omitempty"`
}

type ListDonationsResponseModel struct {
	Donations []Donation `json:"donations"`
}

type CreateDonationRequestModel struct {
	CampaignID    uuid.UUID     `json:"campaign_id"`
	Amount        float64       `json:"amount"`
	DonorID       uuid.UUID     `json:"donor_id"`
	Message       *string       `json:"message,omitempty"`
	IsAnonymous   bool          `json:"is_anonymous"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type CreateDonationResponseModel struct {
	ID uuid.UUID `json:"id"`
}

func MakeEndpointGetDonation(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetDonationRequestModel)
		if !ok {
			return nil, errors.New("failed to cast request")
		}

		donation, err := s.GetDonation(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get donation: %w", err)
		}

		return GetDonationResponseModel{Donation: donation}, nil
	}
}

func MakeEndpointListDonations(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(ListDonationsRequestModel)
		if !ok {
			return nil, errors.New("failed to cast request")
		}

		var donations []Donation

		if req.CampaignID != nil {
			donations, err = s.ListDonationsByCampaign(ctx, *req.CampaignID)
		} else {
			donations, err = s.ListDonations(ctx)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list donations: %w", err)
		}

		return ListDonationsResponseModel{Donations: donations}, nil
	}
}

func MakeEndpointCreateDonation(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateDonationRequestModel)
		if !ok {
			return nil, errors.New("failed to cast request")
		}

		// Validaciones básicas
		if req.Amount <= 0 {
			return nil, apierrors.NewFieldValidationError("amount", "amount must be greater than 0")
		}

		donation := Donation{
			CampaignID:    req.CampaignID,
			Amount:        req.Amount,
			DonorID:       req.DonorID,
			Date:          time.Now(),
			Message:       req.Message,
			IsAnonymous:   req.IsAnonymous,
			PaymentMethod: req.PaymentMethod,
			Status:        DonationStatusPending,
		}

		id, err := s.CreateDonation(ctx, donation)
		if err != nil {
			if _, ok := err.(apierrors.ValidationError); ok {
				return nil, err
			}
			return nil, fmt.Errorf("failed to create donation: %w", err)
		}

		return CreateDonationResponseModel{ID: id}, nil
	}
}
