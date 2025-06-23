package user

import (
	"context"
	"encoding/json"
	"net/http"

	apierrors "dona_tutti_api/errors"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func RegisterRoutes(router *httprouter.Router, s Service) {
	// Options handling for CORS
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

	// Register endpoints
	registerHandler := kithttp.NewServer(
		MakeEndpointRegister(s),
		decodeRegisterRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	loginHandler := kithttp.NewServer(
		MakeEndpointLogin(s),
		decodeLoginRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	getUserHandler := kithttp.NewServer(
		MakeEndpointGetUser(s),
		decodeGetUserRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	updateUserHandler := kithttp.NewServer(
		MakeEndpointUpdateUser(s),
		decodeUpdateUserRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	updatePasswordHandler := kithttp.NewServer(
		MakeEndpointUpdatePassword(s),
		decodeUpdatePasswordRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	requestPasswordResetHandler := kithttp.NewServer(
		MakeEndpointRequestPasswordReset(s),
		decodeRequestPasswordResetRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	resetPasswordHandler := kithttp.NewServer(
		MakeEndpointResetPassword(s),
		decodeResetPasswordRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(apierrors.HTTPErrorEncoder),
	)

	// Register routes
	router.Handle(http.MethodPost, "/auth/register", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		registerHandler.ServeHTTP(w, r)
	}))

	router.Handle(http.MethodPost, "/auth/login", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		loginHandler.ServeHTTP(w, r)
	}))

	router.Handle(http.MethodGet, "/users/:id", withCORS(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		getUserHandler.ServeHTTP(w, r.WithContext(ctx))
	}))

	router.Handle(http.MethodPut, "/users/:id", withCORS(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		updateUserHandler.ServeHTTP(w, r.WithContext(ctx))
	}))

	router.Handle(http.MethodPut, "/users/:id/password", withCORS(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", p)
		updatePasswordHandler.ServeHTTP(w, r.WithContext(ctx))
	}))

	router.Handle(http.MethodPost, "/auth/password-reset", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		requestPasswordResetHandler.ServeHTTP(w, r)
	}))

	router.Handle(http.MethodPost, "/auth/reset-password", withCORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		resetPasswordHandler.ServeHTTP(w, r)
	}))
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RegisterRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request format")
	}
	return req, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request format")
	}
	return req, nil
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, apierrors.NewValidationError("missing route parameters")
	}

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, apierrors.NewFieldValidationError("id", "invalid user ID format")
	}

	return GetUserRequestModel{ID: id}, nil
}

func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, apierrors.NewValidationError("missing route parameters")
	}

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, apierrors.NewFieldValidationError("id", "invalid user ID format")
	}

	var req UpdateUserRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request format")
	}
	req.ID = id
	return req, nil
}

func decodeUpdatePasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params, ok := ctx.Value("params").(httprouter.Params)
	if !ok {
		return nil, apierrors.NewValidationError("missing route parameters")
	}

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, apierrors.NewFieldValidationError("id", "invalid user ID format")
	}

	var req UpdatePasswordRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request format")
	}
	req.ID = id
	return req, nil
}

func decodeRequestPasswordResetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RequestPasswordResetRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request format")
	}
	return req, nil
}

func decodeResetPasswordRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ResetPasswordRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierrors.NewValidationError("invalid request format")
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"data": response,
	})
}
