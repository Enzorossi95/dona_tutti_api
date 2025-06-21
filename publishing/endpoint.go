package publishing

import (
	"context"
	"errors"
	"fmt"
	"microservice_go/article"

	"github.com/go-kit/kit/endpoint"
)

type Service interface {
	GetArticle(ctx context.Context, id string) (article.Article, error)
	CreateArticle(ctx context.Context, article article.Article) (id string, err error)
	ListArticles(ctx context.Context) ([]article.Article, error)
	UpdateArticle(ctx context.Context, id string, article article.Article) (article.Article, error)
}

type GetArticleRequestModel struct {
	ID string `json:"id"`
}

type GetArticleResponseModel struct {
	Article article.Article
}

type ListArticlesRequestModel struct {
	// No parameters needed for listing all articles
}

type ListArticlesResponseModel struct {
	Articles []article.Article
}

type UpdateArticleRequestModel struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateArticleResponseModel struct {
	Article article.Article
}

func MakeEndpointGetArticle(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetArticleRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointGetArticle failed cast request")
		}

		a, err := s.GetArticle(ctx, req.ID)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetArticle: %w", err)
		}

		return GetArticleResponseModel{
			Article: a,
		}, nil
	}
}

func MakeEndpointListArticles(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, ok := request.(ListArticlesRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointListArticles failed cast request")
		}

		articles, err := s.ListArticles(ctx)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointListArticles: %w", err)
		}

		return ListArticlesResponseModel{
			Articles: articles,
		}, nil
	}
}

func MakeEndpointUpdateArticle(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateArticleRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointUpdateArticle failed cast request")
		}

		updatedArticle := article.Article{
			Title:   req.Title,
			Content: req.Content,
		}

		a, err := s.UpdateArticle(ctx, req.ID, updatedArticle)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointUpdateArticle: %w", err)
		}

		return UpdateArticleResponseModel{
			Article: a,
		}, nil
	}
}
