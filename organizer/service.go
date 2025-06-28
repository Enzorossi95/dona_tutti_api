package organizer

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	ListOrganizers(ctx context.Context, userID *uuid.UUID) ([]Organizer, error)
	GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error)
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
