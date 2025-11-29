package organizer

import (
	"context"

	"github.com/google/uuid"
)

// OrganizerInfo represents minimal organizer information for external packages
type OrganizerInfo struct {
	ID      uuid.UUID
	Name    string
	Email   string
	Phone   string
	Address string
}

type Service interface {
	ListOrganizers(ctx context.Context, userID *uuid.UUID) ([]Organizer, error)
	GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error)
	GetOrganizerName(ctx context.Context, organizerID uuid.UUID) (string, error)
	GetOrganizerInfo(ctx context.Context, organizerID uuid.UUID) (OrganizerInfo, error)
	CreateOrganizer(ctx context.Context, organizer Organizer) (Organizer, error)
	UpdateOrganizer(ctx context.Context, organizer Organizer) error
}

type OrganizerRepository interface {
	ListOrganizers(ctx context.Context, userID *uuid.UUID) ([]Organizer, error)
	GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error)
	CreateOrganizer(ctx context.Context, organizer Organizer) (Organizer, error)
	UpdateOrganizer(ctx context.Context, organizer Organizer) error
}

type service struct {
	repo OrganizerRepository
}

func NewService(repo OrganizerRepository) Service {
	return &service{repo: repo}
}

func (s *service) ListOrganizers(ctx context.Context, userID *uuid.UUID) ([]Organizer, error) {
	return s.repo.ListOrganizers(ctx, userID)
}

func (s *service) GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error) {
	return s.repo.GetOrganizer(ctx, id)
}

func (s *service) CreateOrganizer(ctx context.Context, organizer Organizer) (Organizer, error) {
	return s.repo.CreateOrganizer(ctx, organizer)
}

func (s *service) UpdateOrganizer(ctx context.Context, organizer Organizer) error {
	return s.repo.UpdateOrganizer(ctx, organizer)
}

func (s *service) GetOrganizerName(ctx context.Context, organizerID uuid.UUID) (string, error) {
	organizer, err := s.repo.GetOrganizer(ctx, organizerID)
	if err != nil {
		return "", err
	}
	return organizer.Name, nil
}

func (s *service) GetOrganizerInfo(ctx context.Context, organizerID uuid.UUID) (OrganizerInfo, error) {
	organizer, err := s.repo.GetOrganizer(ctx, organizerID)
	if err != nil {
		return OrganizerInfo{}, err
	}
	
	return OrganizerInfo{
		ID:      organizer.ID,
		Name:    organizer.Name,
		Email:   organizer.Email,
		Phone:   organizer.Phone,
		Address: organizer.Address,
	}, nil
}
