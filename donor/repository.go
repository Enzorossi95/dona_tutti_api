package donor

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DonorRepository interface {
	GetDonor(ctx context.Context, id uuid.UUID) (Donor, error)
	CreateDonor(ctx context.Context, donor Donor) error
	UpdateDonor(ctx context.Context, donor Donor) error
	ListDonors(ctx context.Context) ([]Donor, error)
}

type donorRepository struct {
	db *gorm.DB
}

func NewDonorRepository(db *gorm.DB) DonorRepository {
	return &donorRepository{db: db}
}

func (r *donorRepository) GetDonor(ctx context.Context, id uuid.UUID) (Donor, error) {
	var model DonorModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return Donor{}, fmt.Errorf("failed to get donor: %w", err)
	}
	return model.ToEntity(), nil
}

func (r *donorRepository) CreateDonor(ctx context.Context, donor Donor) error {
	model := DonorModel{}
	model.FromEntity(donor)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create donor: %w", err)
	}
	return nil
}

func (r *donorRepository) UpdateDonor(ctx context.Context, donor Donor) error {
	model := DonorModel{}
	model.FromEntity(donor)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("failed to update donor: %w", err)
	}
	return nil
}

func (r *donorRepository) ListDonors(ctx context.Context) ([]Donor, error) {
	var models []DonorModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list donors: %w", err)
	}
	donors := make([]Donor, len(models))
	for i, model := range models {
		donors[i] = model.ToEntity()
	}
	return donors, nil
}
