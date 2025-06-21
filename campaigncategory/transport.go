package campaigncategory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func RegisterRoutes(router *httprouter.Router, s Service) {
	listCategoriesHandler := kithttp.NewServer(
		MakeEndpointListCategories(s),
		decodeListCategoriesRequest,
		encodeListCategoriesResponse,
	)

	getCategoryHandler := kithttp.NewServer(
		MakeEndpointGetCategory(s),
		decodeGetCategoryRequest,
		encodeGetCategoryResponse,
	)

	router.Handle(http.MethodGet, "/categories", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		listCategoriesHandler.ServeHTTP(w, r)
	})

	router.Handle(http.MethodGet, "/categories/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		getCategoryHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decodeListCategoriesRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return ListCategoriesRequestModel{}, nil
}

func decodeGetCategoryRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, fmt.Errorf("failed to get params from context")
	}

	idStr := params.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID format: %w", err)
	}

	return GetCategoryRequestModel{ID: id}, nil
}

func encodeListCategoriesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListCategoriesResponseModel)
	if !ok {
		return fmt.Errorf("encodeListCategoriesResponse failed cast response")
	}
	formatted := formatListCategoriesResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func encodeGetCategoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetCategoryResponseModel)
	if !ok {
		return fmt.Errorf("encodeGetCategoryResponse failed cast response")
	}
	formatted := formatGetCategoryResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func formatListCategoriesResponse(res ListCategoriesResponseModel) map[string]interface{} {
	categories := make([]map[string]interface{}, len(res.Categories))
	for i, category := range res.Categories {
		categories[i] = formatCategory(category)
	}
	return map[string]interface{}{
		"data": map[string]interface{}{
			"categories": categories,
		},
	}
}

func formatGetCategoryResponse(res GetCategoryResponseModel) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"category": formatCategory(res.Category),
		},
	}
}

func formatCategory(category CampaignCategory) map[string]interface{} {
	return map[string]interface{}{
		"id":          category.ID.String(),
		"created_at":  category.CreatedAt.Format(time.RFC3339),
		"name":        category.Name,
		"description": category.Description,
	}
}
