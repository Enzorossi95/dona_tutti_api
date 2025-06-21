package repository

import (
	"context"
	"fmt"
	"microservice_go/article"
	"time"

	"gorm.io/gorm"
)

type articlesRepository struct {
	db *gorm.DB
}

func NewArticlesRepository(db *gorm.DB) *articlesRepository {
	return &articlesRepository{
		db: db,
	}
}

func (r *articlesRepository) GetArticle(ctx context.Context, id string) (article.Article, error) {
	var articleModel article.ArticleModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&articleModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return article.Article{}, fmt.Errorf("article wasn't found")
		}
		return article.Article{}, fmt.Errorf("failed to get article: %w", err)
	}

	return articleModel.ToEntity(), nil
}

func (r *articlesRepository) InsertArticle(ctx context.Context, a article.Article) error {
	var articleModel article.ArticleModel
	articleModel.FromEntity(a)

	err := r.db.WithContext(ctx).Create(&articleModel).Error
	if err != nil {
		return fmt.Errorf("failed to insert article: %w", err)
	}

	return nil
}

func (r *articlesRepository) ListArticles(ctx context.Context) ([]article.Article, error) {
	var articleModels []article.ArticleModel

	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&articleModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list articles: %w", err)
	}

	// Convert to domain entities
	articles := make([]article.Article, len(articleModels))
	for i, model := range articleModels {
		articles[i] = model.ToEntity()
	}

	return articles, nil
}

func (r *articlesRepository) UpdateArticle(ctx context.Context, id string, updatedArticle article.Article) error {
	updatedArticle.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Model(&article.ArticleModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":      updatedArticle.Title,
		"content":    updatedArticle.Content,
		"updated_at": updatedArticle.UpdatedAt,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article with id %s not found", id)
	}

	return nil
}
