package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AlexBBBril/gokit/internal/app/gokit/service"
	"github.com/AlexBBBril/gokit/internal/app/gokit/transport/order"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

// NewService wires Go kit endpoints to the HTTP transport.
func NewService(endpoints order.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
	}

	r.Methods("POST").Path("/orders").Handler(
		kithttp.NewServer(
			endpoints.Create,
			decodeCreateRequest,
			encodeResponse,
			options...,
		))

	// HTTP Post - /orders/{id}
	r.Methods("GET").Path("/orders/{id}").Handler(
		kithttp.NewServer(
			endpoints.GetByID,
			decodeGetByIDRequest,
			encodeResponse,
			options...,
		))

	// HTTP Post - /orders/status
	r.Methods("POST").Path("/orders/status").Handler(
		kithttp.NewServer(
			endpoints.ChangeStatus,
			decodeChangeStatusRequest,
			encodeResponse,
			options...,
		))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req order.CreateRequest
	if err = json.NewDecoder(r.Body).Decode(&req.Order); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return order.GetByIDRequest{ID: id}, nil
}

func decodeChangeStatusRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req order.ChangeStatusRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}

	return req, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)

		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case service.ErrOrderNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
