package campaign

import (
	"context"
	"errors"
	"fmt"
	"time"

	apierrors "microservice_go/errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Service interface {
	GetCampaign(ctx context.Context, id uuid.UUID) (Campaign, error)
	ListCampaigns(ctx context.Context) ([]Campaign, error)
	CreateCampaign(ctx context.Context, campaign Campaign) (uuid.UUID, error)
	GetSummary(ctx context.Context) (Summary, error)
}

type GetCampaignRequestModel struct {
	ID uuid.UUID `json:"id"`
}

type GetCampaignResponseModel struct {
	Campaign Campaign `json:"campaign"`
}

type ListCampaignsRequestModel struct {
	// No parameters needed for listing all campaigns
}

type ListCampaignsResponseModel struct {
	Campaigns []Campaign `json:"campaigns"`
}

type SummaryResponseModel struct {
	Summary Summary `json:"summary"`
}

type CreateCampaignRequestModel struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Goal        float64   `json:"goal"`
	StartDate   string    `json:"start_date"` // Will be parsed to time.Time
	EndDate     string    `json:"end_date"`   // Will be parsed to time.Time
	Location    string    `json:"location"`
	CategoryId  uuid.UUID `json:"category"`
	Urgency     int       `json:"urgency"`
	OrganizerId uuid.UUID `json:"organizer"`
}

type CreateCampaignResponseModel struct {
	ID uuid.UUID `json:"id"`
}

func MakeEndpointGetCampaign(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetCampaignRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointGetCampaign failed cast request")
		}

		campaign, err := s.GetCampaign(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetCampaign: %w", err)
		}

		return GetCampaignResponseModel{
			Campaign: campaign,
		}, nil
	}
}

func MakeEndpointListCampaigns(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, ok := request.(ListCampaignsRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointListCampaigns failed cast request")
		}

		campaigns, err := s.ListCampaigns(ctx)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointListCampaigns: %w", err)
		}

		return ListCampaignsResponseModel{
			Campaigns: campaigns,
		}, nil
	}
}

func MakeEndpointCreateCampaign(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateCampaignRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointCreateCampaign failed cast request")
		}

		campaign := Campaign{
			Title:       req.Title,
			Description: req.Description,
			Image:       req.Image,
			Goal:        req.Goal,
			Location:    req.Location,
			CategoryId:  req.CategoryId,
			Urgency:     req.Urgency,
			OrganizerId: req.OrganizerId,
		}

		// Parse dates if provided
		if req.StartDate != "" {
			startDate, err := time.Parse(time.RFC3339, req.StartDate)
			if err != nil {
				return nil, apierrors.NewFieldValidationError("start_date", "invalid start_date format, use RFC3339")
			}
			campaign.StartDate = startDate
		}

		if req.EndDate != "" {
			endDate, err := time.Parse(time.RFC3339, req.EndDate)
			if err != nil {
				return nil, apierrors.NewFieldValidationError("end_date", "invalid end_date format, use RFC3339")
			}
			campaign.EndDate = endDate
		}

		id, err := s.CreateCampaign(ctx, campaign)
		if err != nil {
			// Check if it's a validation error and return it directly to preserve the type
			if _, ok := err.(apierrors.ValidationError); ok {
				return nil, err
			}
			// For other errors, wrap them
			return nil, fmt.Errorf("MakeEndpointCreateCampaign: %w", err)
		}

		return CreateCampaignResponseModel{
			ID: id,
		}, nil
	}
}

func MakeEndpointGetSummary(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		summary, err := s.GetSummary(ctx)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetSummary: %w", err)
		}

		return SummaryResponseModel{
			Summary: summary,
		}, nil
	}
}
