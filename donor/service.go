package donor

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GetDonor(ctx context.Context, id uuid.UUID) (Donor, error)
	CreateDonor(ctx context.Context, donor Donor) (uuid.UUID, error)
	UpdateDonor(ctx context.Context, donor Donor) error
	ListDonors(ctx context.Context) ([]Donor, error)
	FindDonorByEmail(ctx context.Context, email string) (Donor, error)
	FindDonorByPhone(ctx context.Context, phone string) (Donor, error)
}

type service struct {
	repo DonorRepository
}

func NewService(repo DonorRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetDonor(ctx context.Context, id uuid.UUID) (Donor, error) {
	return s.repo.GetDonor(ctx, id)
}

func (s *service) CreateDonor(ctx context.Context, donor Donor) (uuid.UUID, error) {
	donor.ID = uuid.New()
	if err := s.repo.CreateDonor(ctx, donor); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create donor: %w", err)
	}
	return donor.ID, nil
}

func (s *service) UpdateDonor(ctx context.Context, donor Donor) error {
	return s.repo.UpdateDonor(ctx, donor)
}

func (s *service) ListDonors(ctx context.Context) ([]Donor, error) {
	return s.repo.ListDonors(ctx)
}

func (s *service) FindDonorByEmail(ctx context.Context, email string) (Donor, error) {
	return s.repo.FindDonorByEmail(ctx, email)
}

func (s *service) FindDonorByPhone(ctx context.Context, phone string) (Donor, error) {
	return s.repo.FindDonorByPhone(ctx, phone)
}
