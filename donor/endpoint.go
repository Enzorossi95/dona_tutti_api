package donor

import (
	"context"
	"errors"
	"fmt"

	apierrors "dona_tutti_api/errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type GetDonorRequestModel struct {
	ID uuid.UUID `json:"id"`
}

type GetDonorResponseModel struct {
	Donor Donor `json:"donor"`
}

type ListDonorsRequestModel struct {
	// No parameters needed
}

type ListDonorsResponseModel struct {
	Donors []Donor `json:"donors"`
}

type CreateDonorRequestModel struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type CreateDonorResponseModel struct {
	ID uuid.UUID `json:"id"`
}

func MakeEndpointGetDonor(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetDonorRequestModel)
		if !ok {
			return nil, errors.New("failed to cast request")
		}

		donor, err := s.GetDonor(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get donor: %w", err)
		}

		return GetDonorResponseModel{Donor: donor}, nil
	}
}

func MakeEndpointListDonors(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, ok := request.(ListDonorsRequestModel)
		if !ok {
			return nil, errors.New("failed to cast request")
		}

		donors, err := s.ListDonors(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list donors: %w", err)
		}

		return ListDonorsResponseModel{Donors: donors}, nil
	}
}

func MakeEndpointCreateDonor(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateDonorRequestModel)
		if !ok {
			return nil, errors.New("failed to cast request")
		}

		donor := Donor{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Phone:     req.Phone,
			Email:     req.Email,
		}

		id, err := s.CreateDonor(ctx, donor)
		if err != nil {
			if _, ok := err.(apierrors.ValidationError); ok {
				return nil, err
			}
			return nil, fmt.Errorf("failed to create donor: %w", err)
		}

		return CreateDonorResponseModel{ID: id}, nil
	}
}
