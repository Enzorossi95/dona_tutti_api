package organizer

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
	listOrganizersHandler := kithttp.NewServer(
		MakeEndpointListOrganizers(s),
		decodeListOrganizersRequest,
		encodeListOrganizersResponse,
	)

	getOrganizerHandler := kithttp.NewServer(
		MakeEndpointGetOrganizer(s),
		decodeGetOrganizerRequest,
		encodeGetOrganizerResponse,
	)

	router.Handle(http.MethodGet, "/organizers", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		listOrganizersHandler.ServeHTTP(w, r)
	})

	router.Handle(http.MethodGet, "/organizers/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		getOrganizerHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decodeListOrganizersRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return ListOrganizersRequestModel{}, nil
}

func decodeGetOrganizerRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, fmt.Errorf("failed to get params from context")
	}

	idStr := params.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid organizer ID format: %w", err)
	}

	return GetOrganizerRequestModel{ID: id}, nil
}

func encodeListOrganizersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListOrganizersResponseModel)
	if !ok {
		return fmt.Errorf("encodeListOrganizersResponse failed cast response")
	}
	formatted := formatListOrganizersResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func encodeGetOrganizerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetOrganizerResponseModel)
	if !ok {
		return fmt.Errorf("encodeGetOrganizerResponse failed cast response")
	}
	formatted := formatGetOrganizerResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func formatListOrganizersResponse(res ListOrganizersResponseModel) map[string]interface{} {
	organizers := make([]map[string]interface{}, len(res.Organizers))
	for i, organizer := range res.Organizers {
		organizers[i] = formatOrganizer(organizer)
	}
	return map[string]interface{}{
		"data": map[string]interface{}{
			"organizers": organizers,
		},
	}
}

func formatGetOrganizerResponse(res GetOrganizerResponseModel) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"organizer": formatOrganizer(res.Organizer),
		},
	}
}

func formatOrganizer(organizer Organizer) map[string]interface{} {
	return map[string]interface{}{
		"id":         organizer.ID.String(),
		"created_at": organizer.CreatedAt.Format(time.RFC3339),
		"name":       organizer.Name,
		"avatar":     organizer.Avatar,
		"verified":   organizer.Verified,
	}
}
