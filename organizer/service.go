package organizer

import (
	"context"

	"github.com/google/uuid"
)

type OrganizerRepository interface {
	ListOrganizers(ctx context.Context) ([]Organizer, error)
	GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error)
}

type service struct {
	repo OrganizerRepository
}

func NewService(repo OrganizerRepository) *service {
	return &service{repo: repo}
}

func (s *service) ListOrganizers(ctx context.Context) ([]Organizer, error) {
	return s.repo.ListOrganizers(ctx)
}

func (s *service) GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error) {
	return s.repo.GetOrganizer(ctx, id)
}
