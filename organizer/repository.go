package organizer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type organizerRepository struct {
	db *gorm.DB
}

func NewOrganizerRepository(db *gorm.DB) *organizerRepository {
	return &organizerRepository{
		db: db,
	}
}

func (r *organizerRepository) ListOrganizers(ctx context.Context) ([]Organizer, error) {
	var organizerModels []OrganizerModel

	err := r.db.WithContext(ctx).Order("name").Find(&organizerModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list organizers: %w", err)
	}

	// Convert to domain entities
	organizers := make([]Organizer, len(organizerModels))
	for i, model := range organizerModels {
		organizers[i] = model.ToEntity()
	}

	return organizers, nil
}

func (r *organizerRepository) GetOrganizer(ctx context.Context, id uuid.UUID) (Organizer, error) {
	var organizerModel OrganizerModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&organizerModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Organizer{}, fmt.Errorf("organizer with id %s not found", id.String())
		}
		return Organizer{}, fmt.Errorf("failed to get organizer: %w", err)
	}

	return organizerModel.ToEntity(), nil
}
