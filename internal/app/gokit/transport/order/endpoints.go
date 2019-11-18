package order

import (
	"context"
	"github.com/AlexBBBril/gokit/internal/app/gokit/service"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	ChangeStatus endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s service.OrderService) Endpoints {
	return Endpoints{
		Create:       makeCreateEndpoint(s),
		GetByID:      makeGetByIDEndpoint(s),
		ChangeStatus: makeChangeStatusEndpoint(s),
	}
}

func makeCreateEndpoint(s service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Order)

		return CreateResponse{
			ID:  id,
			Err: err,
		}, nil
	}
}

func makeGetByIDEndpoint(s service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetByIDRequest)
		orderRequest, err := s.GetByID(ctx, req.ID)

		return GetByIDResponse{
			Order: orderRequest,
			Err:   err,
		}, nil
	}
}

func makeChangeStatusEndpoint(s service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ChangeStatusRequest)
		err = s.ChangeStatus(ctx, req.ID, req.Status)

		return ChangeStatusResponse{Err: err}, nil
	}
}
