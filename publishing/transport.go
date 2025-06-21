package publishing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

type Router interface {
	Handle(method, path string, handler http.Handler)
}

func RegisterRoutes(router *httprouter.Router, s Service) {
	getArticleHandler := kithttp.NewServer(
		MakeEndpointGetArticle(s),
		decodeGetArticleRequest,
		encodeGetArticleResponse,
	)

	listArticlesHandler := kithttp.NewServer(
		MakeEndpointListArticles(s),
		decodeListArticlesRequest,
		encodeListArticlesResponse,
	)

	updateArticleHandler := kithttp.NewServer(
		MakeEndpointUpdateArticle(s),
		decodeUpdateArticleRequest,
		encodeUpdateArticleResponse,
	)

	router.Handle(http.MethodGet, "/articles/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Agregamos los parámetros de la URL al contexto para que estén disponibles en decodeGetArticleRequest
		ctx := context.WithValue(r.Context(), "params", p)
		getArticleHandler.ServeHTTP(w, r.WithContext(ctx))
	})

	router.Handle(http.MethodGet, "/articles", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		listArticlesHandler.ServeHTTP(w, r)
	})

	router.Handle(http.MethodPut, "/articles/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Agregamos los parámetros de la URL al contexto para que estén disponibles en decodeUpdateArticleRequest
		ctx := context.WithValue(r.Context(), "params", p)
		updateArticleHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decodeGetArticleRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, fmt.Errorf("failed to get params from context")
	}
	return GetArticleRequestModel{ID: params.ByName("id")}, nil
}

func encodeGetArticleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetArticleResponseModel)
	if !ok {
		return fmt.Errorf("encodeGetArticleResponse failed cast response")
	}
	formatted := formatGetArticleResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func formatGetArticleResponse(res GetArticleResponseModel) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"article": map[string]interface{}{
				"id":      res.Article.ID,
				"title":   res.Article.Title,
				"content": res.Article.Content,
			},
		},
	}
}

func decodeListArticlesRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return ListArticlesRequestModel{}, nil
}

func encodeListArticlesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListArticlesResponseModel)
	if !ok {
		return fmt.Errorf("encodeListArticlesResponse failed cast response")
	}
	formatted := formatListArticlesResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func formatListArticlesResponse(res ListArticlesResponseModel) map[string]interface{} {
	articles := make([]map[string]interface{}, len(res.Articles))
	for i, article := range res.Articles {
		articles[i] = map[string]interface{}{
			"id":      article.ID,
			"title":   article.Title,
			"content": article.Content,
		}
	}
	return map[string]interface{}{
		"data": map[string]interface{}{
			"articles": articles,
		},
	}
}

func decodeUpdateArticleRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, fmt.Errorf("failed to get params from context")
	}

	var req UpdateArticleRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	// Set the ID from URL parameter
	req.ID = params.ByName("id")

	return req, nil
}

func encodeUpdateArticleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(UpdateArticleResponseModel)
	if !ok {
		return fmt.Errorf("encodeUpdateArticleResponse failed cast response")
	}
	formatted := formatUpdateArticleResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func formatUpdateArticleResponse(res UpdateArticleResponseModel) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"article": map[string]interface{}{
				"id":      res.Article.ID,
				"title":   res.Article.Title,
				"content": res.Article.Content,
			},
		},
	}
}
