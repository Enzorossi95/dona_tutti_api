package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type RegisterRequestModel struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterResponseModel struct {
	ID uuid.UUID `json:"id"`
}

type LoginRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseModel struct {
	Token AuthToken `json:"token"`
}

type GetUserRequestModel struct {
	ID uuid.UUID `json:"id"`
}

type GetUserResponseModel struct {
	User User `json:"user"`
}

type UpdateUserRequestModel struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type UpdatePasswordRequestModel struct {
	ID              uuid.UUID `json:"id"`
	CurrentPassword string    `json:"current_password"`
	NewPassword     string    `json:"new_password"`
}

type RequestPasswordResetRequestModel struct {
	Email string `json:"email"`
}

type ResetPasswordRequestModel struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func MakeEndpointRegister(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(RegisterRequestModel)
		if !ok {
			return nil, errors.New("invalid request format")
		}

		id, err := s.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName)
		if err != nil {
			return nil, fmt.Errorf("failed to register user: %w", err)
		}

		return RegisterResponseModel{ID: id}, nil
	}
}

func MakeEndpointLogin(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(LoginRequestModel)
		if !ok {
			return nil, errors.New("invalid request format")
		}

		token, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to login: %w", err)
		}

		return LoginResponseModel{Token: *token}, nil
	}
}

func MakeEndpointGetUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetUserRequestModel)
		if !ok {
			return nil, errors.New("invalid request format")
		}

		user, err := s.GetUser(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		return GetUserResponseModel{User: user}, nil
	}
}

func MakeEndpointUpdateUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateUserRequestModel)
		if !ok {
			return nil, errors.New("invalid request format")
		}

		err = s.UpdateUser(ctx, req.ID, req.FirstName, req.LastName)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		return struct{}{}, nil
	}
}

func MakeEndpointUpdatePassword(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdatePasswordRequestModel)
		if !ok {
			return nil, errors.New("invalid request format")
		}

		err = s.UpdatePassword(ctx, req.ID, req.CurrentPassword, req.NewPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to update password: %w", err)
		}

		return struct{}{}, nil
	}
}

func MakeEndpointRequestPasswordReset(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(RequestPasswordResetRequestModel)
		if !ok {
			return nil, errors.New("invalid request format")
		}

		err = s.RequestPasswordReset(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to request password reset: %w", err)
		}

		return struct{}{}, nil
	}
}

func MakeEndpointResetPassword(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(ResetPasswordRequestModel)
		if !ok {
			return nil, errors.New("invalid request format")
		}

		err = s.ResetPassword(ctx, req.Token, req.NewPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to reset password: %w", err)
		}

		return struct{}{}, nil
	}
}
