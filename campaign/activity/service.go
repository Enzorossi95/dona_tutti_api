package activity

import (
	"context"
	"fmt"
	"time"

	apierrors "dona_tutti_api/errors"

	"github.com/google/uuid"
)

type Service interface {
	GetActivitiesByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Activity, error)
	GetActivity(ctx context.Context, id uuid.UUID) (Activity, error)
	CreateActivity(ctx context.Context, activity Activity) (uuid.UUID, error)
	UpdateActivity(ctx context.Context, id uuid.UUID, campaignID uuid.UUID, updateReq ActivityUpdateRequest) error
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type ActivityRepository interface {
	GetActivitiesByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Activity, error)
	GetActivity(ctx context.Context, id uuid.UUID) (Activity, error)
	CreateActivity(ctx context.Context, activity Activity) error
	UpdateActivity(ctx context.Context, activity Activity) error
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo ActivityRepository
}

func NewService(repo ActivityRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetActivitiesByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Activity, error) {
	return s.repo.GetActivitiesByCampaign(ctx, campaignID)
}

func (s *service) GetActivity(ctx context.Context, id uuid.UUID) (Activity, error) {
	return s.repo.GetActivity(ctx, id)
}

func (s *service) CreateActivity(ctx context.Context, activity Activity) (uuid.UUID, error) {
	// Generate new ID and set timestamps
	activity.ID = uuid.New()
	activity.CreatedAt = time.Now()

	// Validate required fields
	if activity.CampaignID == uuid.Nil {
		return uuid.Nil, apierrors.NewFieldValidationError("campaign_id", "campaign ID is required")
	}
	if activity.Title == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("title", "activity title is required")
	}
	if activity.Type == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("type", "activity type is required")
	}
	if activity.Author == "" {
		return uuid.Nil, apierrors.NewFieldValidationError("author", "activity author is required")
	}
	if activity.Date.IsZero() {
		return uuid.Nil, apierrors.NewFieldValidationError("date", "activity date is required")
	}

	if err := s.repo.CreateActivity(ctx, activity); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return activity.ID, nil
}

func (s *service) UpdateActivity(ctx context.Context, id uuid.UUID, campaignID uuid.UUID, updateReq ActivityUpdateRequest) error {
	// Get existing activity
	existingActivity, err := s.repo.GetActivity(ctx, id)
	if err != nil {
		return fmt.Errorf("activity not found: %w", err)
	}

	// Verify that the activity belongs to the specified campaign
	if existingActivity.CampaignID != campaignID {
		return apierrors.NewFieldValidationError("campaign_id", "activity does not belong to the specified campaign")
	}

	// Apply partial updates - only update fields that are provided
	updatedActivity := existingActivity

	if updateReq.Title != nil {
		if *updateReq.Title == "" {
			return apierrors.NewFieldValidationError("title", "activity title cannot be empty")
		}
		updatedActivity.Title = *updateReq.Title
	}

	if updateReq.Description != nil {
		updatedActivity.Description = *updateReq.Description
	}

	if updateReq.Date != nil {
		if updateReq.Date.IsZero() {
			return apierrors.NewFieldValidationError("date", "activity date cannot be empty")
		}
		updatedActivity.Date = *updateReq.Date
	}

	if updateReq.Type != nil {
		if *updateReq.Type == "" {
			return apierrors.NewFieldValidationError("type", "activity type cannot be empty")
		}
		updatedActivity.Type = *updateReq.Type
	}

	if updateReq.Author != nil {
		if *updateReq.Author == "" {
			return apierrors.NewFieldValidationError("author", "activity author cannot be empty")
		}
		updatedActivity.Author = *updateReq.Author
	}

	// Update the activity
	if err := s.repo.UpdateActivity(ctx, updatedActivity); err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	return nil
}

func (s *service) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	// Check if activity exists
	if _, err := s.repo.GetActivity(ctx, id); err != nil {
		return fmt.Errorf("activity not found: %w", err)
	}

	if err := s.repo.DeleteActivity(ctx, id); err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	return nil
}