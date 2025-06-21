package campaigncategory

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Service interface {
	ListCategories(ctx context.Context) ([]CampaignCategory, error)
	GetCategory(ctx context.Context, id uuid.UUID) (CampaignCategory, error)
}

type ListCategoriesRequestModel struct{}

type ListCategoriesResponseModel struct {
	Categories []CampaignCategory `json:"categories"`
}

type GetCategoryRequestModel struct {
	ID uuid.UUID `json:"id"`
}

type GetCategoryResponseModel struct {
	Category CampaignCategory `json:"category"`
}

func MakeEndpointListCategories(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, ok := request.(ListCategoriesRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointListCategories failed cast request")
		}

		categories, err := s.ListCategories(ctx)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointListCategories: %w", err)
		}

		return ListCategoriesResponseModel{
			Categories: categories,
		}, nil
	}
}

func MakeEndpointGetCategory(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetCategoryRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointGetCategory failed cast request")
		}

		category, err := s.GetCategory(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetCategory: %w", err)
		}

		return GetCategoryResponseModel{
			Category: category,
		}, nil
	}
}
