package organizer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository implementation
type organizerRepository struct {
	db *gorm.DB
}

// NewOrganizerRepository creates a new organizer repository
func NewOrganizerRepository(db *gorm.DB) OrganizerRepository {
	return &organizerRepository{db: db}
}

func (r *organizerRepository) ListOrganizers(ctx context.Context, userID *uuid.UUID) ([]Organizer, error) {
	var organizerModels []OrganizerModel

	query := r.db.WithContext(ctx).Order("name")

	// Apply user ID filter if provided
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	err := query.Find(&organizerModels).Error
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

func (r *organizerRepository) CreateOrganizer(ctx context.Context, organizer Organizer) (Organizer, error) {
	organizerModel := OrganizerModel{}
	organizerModel.FromEntity(organizer)

	err := r.db.WithContext(ctx).Create(&organizerModel).Error
	if err != nil {
		return Organizer{}, fmt.Errorf("failed to create organizer: %w", err)
	}

	return organizerModel.ToEntity(), nil
}

func (r *organizerRepository) UpdateOrganizer(ctx context.Context, organizer Organizer) error {
	model := OrganizerModel{}
	model.FromEntity(organizer)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("failed to update organizer: %w", err)
	}
	return nil
}
