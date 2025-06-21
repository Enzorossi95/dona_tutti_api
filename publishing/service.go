package publishing

import (
	"context"
	"fmt"
	"math/rand"
	"microservice_go/article"
	"strconv"
	"time"
)

type ArticlesRepository interface {
	GetArticle(ctx context.Context, id string) (article.Article, error)
	InsertArticle(ctx context.Context, article article.Article) error
	ListArticles(ctx context.Context) ([]article.Article, error)
	UpdateArticle(ctx context.Context, id string, article article.Article) error
}

type service struct {
	repo ArticlesRepository
}

func NewService(repo ArticlesRepository) *service {
	return &service{repo: repo}
}

func (s *service) GetArticle(ctx context.Context, id string) (article.Article, error) {
	return s.repo.GetArticle(ctx, id)
}

func (s *service) CreateArticle(ctx context.Context, article article.Article) (id string, err error) {
	article.ID = generateID()

	if err := s.repo.InsertArticle(ctx, article); err != nil {
		return "", fmt.Errorf("failed to insert article: %w", err)
	}

	return article.ID, nil
}

func (s *service) ListArticles(ctx context.Context) ([]article.Article, error) {
	return s.repo.ListArticles(ctx)
}

func (s *service) UpdateArticle(ctx context.Context, id string, updatedArticle article.Article) (article.Article, error) {
	// Set the UpdatedAt timestamp
	updatedArticle.UpdatedAt = time.Now()

	if err := s.repo.UpdateArticle(ctx, id, updatedArticle); err != nil {
		return article.Article{}, fmt.Errorf("failed to update article: %w", err)
	}

	// Return the updated article
	return s.repo.GetArticle(ctx, id)
}

func generateID() string {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)
	return strconv.FormatInt(timestamp, 10) + strconv.Itoa(random)
}
