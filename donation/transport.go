package donation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apierrors "microservice_go/errors"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func RegisterRoutes(router *httprouter.Router, s Service) {
	getDonationHandler := kithttp.NewServer(
		MakeEndpointGetDonation(s),
		decodeGetDonationRequest,
		encodeGetDonationResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	listDonationsHandler := kithttp.NewServer(
		MakeEndpointListDonations(s),
		decodeListDonationsRequest,
		encodeListDonationsResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	createDonationHandler := kithttp.NewServer(
		MakeEndpointCreateDonation(s),
		decodeCreateDonationRequest,
		encodeCreateDonationResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	// Wrapper para CORS
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

	router.Handle(http.MethodGet, "/donations/:id", withCORS(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		getDonationHandler.ServeHTTP(w, r.WithContext(ctx))
	}))

	router.Handle(http.MethodGet, "/donations", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		listDonationsHandler.ServeHTTP(w, r)
	}))

	router.Handle(http.MethodPost, "/donations", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		createDonationHandler.ServeHTTP(w, r)
	}))
}

func decodeGetDonationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, fmt.Errorf("failed to get params from context")
	}

	idStr := params.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, apierrors.NewFieldValidationError("id", "invalid donation ID format")
	}

	return GetDonationRequestModel{ID: id}, nil
}

func decodeListDonationsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	campaignIDStr := r.URL.Query().Get("campaign_id")
	if campaignIDStr == "" {
		return ListDonationsRequestModel{}, nil
	}

	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return nil, apierrors.NewFieldValidationError("campaign_id", "invalid campaign ID format")
	}

	return ListDonationsRequestModel{CampaignID: &campaignID}, nil
}

func decodeCreateDonationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateDonationRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request body")
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeGetDonationResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetDonationResponseModel)
	if !ok {
		return fmt.Errorf("failed to cast response")
	}
	return encodeResponse(ctx, w, map[string]interface{}{
		"data": map[string]interface{}{
			"donation": res.Donation,
		},
	})
}

func encodeListDonationsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListDonationsResponseModel)
	if !ok {
		return fmt.Errorf("failed to cast response")
	}
	return encodeResponse(ctx, w, map[string]interface{}{
		"data": map[string]interface{}{
			"donations": res.Donations,
		},
	})
}

func encodeCreateDonationResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(CreateDonationResponseModel)
	if !ok {
		return fmt.Errorf("failed to cast response")
	}
	w.WriteHeader(http.StatusCreated)
	return encodeResponse(ctx, w, map[string]interface{}{
		"data": map[string]interface{}{
			"id": res.ID,
		},
	})
}
