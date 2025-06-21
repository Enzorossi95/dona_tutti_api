package campaigncategory

import (
	"context"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]CampaignCategory, error)
	GetCategory(ctx context.Context, id uuid.UUID) (CampaignCategory, error)
}

type service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) *service {
	return &service{repo: repo}
}

func (s *service) ListCategories(ctx context.Context) ([]CampaignCategory, error) {
	return s.repo.ListCategories(ctx)
}

func (s *service) GetCategory(ctx context.Context, id uuid.UUID) (CampaignCategory, error) {
	return s.repo.GetCategory(ctx, id)
}
