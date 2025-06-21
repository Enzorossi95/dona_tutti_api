package campaigncategory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) ListCategories(ctx context.Context) ([]CampaignCategory, error) {
	var categoryModels []CampaignCategoryModel

	err := r.db.WithContext(ctx).Order("name").Find(&categoryModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	// Convert to domain entities
	categories := make([]CampaignCategory, len(categoryModels))
	for i, model := range categoryModels {
		categories[i] = model.ToEntity()
	}

	return categories, nil
}

func (r *categoryRepository) GetCategory(ctx context.Context, id uuid.UUID) (CampaignCategory, error) {
	var categoryModel CampaignCategoryModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&categoryModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return CampaignCategory{}, fmt.Errorf("category with id %s not found", id.String())
		}
		return CampaignCategory{}, fmt.Errorf("failed to get category: %w", err)
	}

	return categoryModel.ToEntity(), nil
}
