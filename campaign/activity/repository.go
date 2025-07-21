package activity

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetActivitiesByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Activity, error)
	GetActivity(ctx context.Context, id uuid.UUID) (Activity, error)
	CreateActivity(ctx context.Context, activity Activity) error
	UpdateActivity(ctx context.Context, activity Activity) error
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetActivitiesByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Activity, error) {
	var models []ActivityModel
	if err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Order("date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	activities := make([]Activity, len(models))
	for i, model := range models {
		activities[i] = model.ToEntity()
	}

	return activities, nil
}

func (r *repository) GetActivity(ctx context.Context, id uuid.UUID) (Activity, error) {
	var model ActivityModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return Activity{}, err
	}

	return model.ToEntity(), nil
}

func (r *repository) CreateActivity(ctx context.Context, activity Activity) error {
	var model ActivityModel
	model.FromEntity(activity)

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *repository) UpdateActivity(ctx context.Context, activity Activity) error {
	var model ActivityModel
	model.FromEntity(activity)

	return r.db.WithContext(ctx).Where("id = ?", activity.ID).Updates(&model).Error
}

func (r *repository) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ActivityModel{}).Error
}