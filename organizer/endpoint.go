package organizer

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Service interface {
	ListOrganizers(ctx context.Context) ([]Organizer, error)
	GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error)
}

type ListOrganizersRequestModel struct{}

type ListOrganizersResponseModel struct {
	Organizers []Organizer `json:"organizers"`
}

type GetOrganizerRequestModel struct {
	ID uuid.UUID `json:"id"`
}

type GetOrganizerResponseModel struct {
	Organizer Organizer `json:"organizer"`
}

func MakeEndpointListOrganizers(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, ok := request.(ListOrganizersRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointListOrganizers failed cast request")
		}

		organizers, err := s.ListOrganizers(ctx)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointListOrganizers: %w", err)
		}

		return ListOrganizersResponseModel{
			Organizers: organizers,
		}, nil
	}
}

func MakeEndpointGetOrganizer(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetOrganizerRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointGetOrganizer failed cast request")
		}

		organizer, err := s.GetOrganizer(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetOrganizer: %w", err)
		}

		return GetOrganizerResponseModel{
			Organizer: organizer,
		}, nil
	}
}
