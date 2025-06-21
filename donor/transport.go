package donor

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
	getDonorHandler := kithttp.NewServer(
		MakeEndpointGetDonor(s),
		decodeGetDonorRequest,
		encodeGetDonorResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	listDonorsHandler := kithttp.NewServer(
		MakeEndpointListDonors(s),
		decodeListDonorsRequest,
		encodeListDonorsResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	createDonorHandler := kithttp.NewServer(
		MakeEndpointCreateDonor(s),
		decodeCreateDonorRequest,
		encodeCreateDonorResponse,
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

	router.Handle(http.MethodGet, "/donors/:id", withCORS(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		getDonorHandler.ServeHTTP(w, r.WithContext(ctx))
	}))

	router.Handle(http.MethodGet, "/donors", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		listDonorsHandler.ServeHTTP(w, r)
	}))

	router.Handle(http.MethodPost, "/donors", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		createDonorHandler.ServeHTTP(w, r)
	}))
}

func decodeGetDonorRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, fmt.Errorf("failed to get params from context")
	}

	idStr := params.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, apierrors.NewFieldValidationError("id", "invalid donor ID format")
	}

	return GetDonorRequestModel{ID: id}, nil
}

func decodeListDonorsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return ListDonorsRequestModel{}, nil
}

func decodeCreateDonorRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateDonorRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request body")
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeGetDonorResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetDonorResponseModel)
	if !ok {
		return fmt.Errorf("failed to cast response")
	}
	return encodeResponse(ctx, w, map[string]interface{}{
		"data": map[string]interface{}{
			"donor": res.Donor,
		},
	})
}

func encodeListDonorsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListDonorsResponseModel)
	if !ok {
		return fmt.Errorf("failed to cast response")
	}
	return encodeResponse(ctx, w, map[string]interface{}{
		"data": map[string]interface{}{
			"donors": res.Donors,
		},
	})
}

func encodeCreateDonorResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(CreateDonorResponseModel)
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
