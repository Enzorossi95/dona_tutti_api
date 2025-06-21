package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apierrors "microservice_go/errors"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// Agregar esta función al inicio del archivo
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // En producción, especifica el origen exacto
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterRoutes(router *httprouter.Router, s Service) {
	getCampaignHandler := kithttp.NewServer(
		MakeEndpointGetCampaign(s),
		decodeGetCampaignRequest,
		encodeGetCampaignResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	listCampaignsHandler := kithttp.NewServer(
		MakeEndpointListCampaigns(s),
		decodeListCampaignsRequest,
		encodeListCampaignsResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	createCampaignHandler := kithttp.NewServer(
		MakeEndpointCreateCampaign(s),
		decodeCreateCampaignRequest,
		encodeCreateCampaignResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	summaryHandler := kithttp.NewServer(
		MakeEndpointGetSummary(s),
		decodeSummaryRequest,
		encodeSummaryResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	// Modificar el wrapper para que maneje correctamente los tipos
	withCORS := func(handle httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			handle(w, r, ps)
		}
	}

	// Usar el wrapper modificado
	router.Handle(http.MethodGet, "/summary/campaigns", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		summaryHandler.ServeHTTP(w, r)
	}))

	router.Handle(http.MethodGet, "/campaigns/:id", withCORS(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		getCampaignHandler.ServeHTTP(w, r.WithContext(ctx))
	}))

	router.Handle(http.MethodGet, "/campaigns", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		listCampaignsHandler.ServeHTTP(w, r)
	}))

	router.Handle(http.MethodPost, "/campaigns", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		createCampaignHandler.ServeHTTP(w, r)
	}))

	// Ya no necesitamos GlobalOPTIONS porque el middleware maneja OPTIONS
}

func decodeGetCampaignRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, fmt.Errorf("failed to get params from context")
	}

	idStr := params.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, apierrors.NewFieldValidationError("id", "invalid campaign ID format")
	}

	return GetCampaignRequestModel{ID: id}, nil
}

func decodeListCampaignsRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return ListCampaignsRequestModel{}, nil
}

func decodeCreateCampaignRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req CreateCampaignRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("failed to decode request body: invalid JSON format")
	}
	return req, nil
}

func encodeGetCampaignResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetCampaignResponseModel)
	if !ok {
		return fmt.Errorf("encodeGetCampaignResponse failed cast response")
	}
	formatted := formatGetCampaignResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func encodeListCampaignsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListCampaignsResponseModel)
	if !ok {
		return fmt.Errorf("encodeListCampaignsResponse failed cast response")
	}
	formatted := formatListCampaignsResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func encodeCreateCampaignResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(CreateCampaignResponseModel)
	if !ok {
		return fmt.Errorf("encodeCreateCampaignResponse failed cast response")
	}
	formatted := formatCreateCampaignResponse(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(formatted)
}

func formatGetCampaignResponse(res GetCampaignResponseModel) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"campaign": formatCampaign(res.Campaign),
		},
	}
}

func formatListCampaignsResponse(res ListCampaignsResponseModel) map[string]interface{} {
	campaigns := make([]map[string]interface{}, len(res.Campaigns))
	for i, campaign := range res.Campaigns {
		campaigns[i] = formatCampaign(campaign)
	}
	return map[string]interface{}{
		"data": map[string]interface{}{
			"campaigns": campaigns,
		},
	}
}

func formatCreateCampaignResponse(res CreateCampaignResponseModel) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"id": res.ID.String(),
		},
	}
}

func formatCampaign(campaign Campaign) map[string]interface{} {
	return map[string]interface{}{
		"id":          campaign.ID.String(),
		"created_at":  campaign.CreatedAt.Format(time.RFC3339),
		"title":       campaign.Title,
		"description": campaign.Description,
		"image":       campaign.Image,
		"goal":        campaign.Goal,
		"start_date":  campaign.StartDate.Format(time.RFC3339),
		"end_date":    campaign.EndDate.Format(time.RFC3339),
		"location":    campaign.Location,
		"category":    campaign.CategoryId,
		"urgency":     campaign.Urgency,
		"organizer":   campaign.OrganizerId,
		"status":      campaign.Status,
	}
}

func formatSummaryResponse(res Summary) map[string]interface{} {
	return map[string]interface{}{
		"result": map[string]interface{}{
			"total_campaigns":    res.TotalCampaigns,
			"total_goal":         res.TotalGoal,
			"total_contributors": res.TotalContributors,
		},
	}
}

func encodeSummaryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(SummaryResponseModel)
	if !ok {
		return fmt.Errorf("encodeSummaryResponse failed cast response")
	}
	formatted := formatSummaryResponse(res.Summary)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}

func decodeSummaryRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return struct{}{}, nil
}
