package endpoint

import (
	"context"
	"github.com/AlexBBBril/gokit/internal/app/gokit/service"
	"github.com/AlexBBBril/gokit/internal/app/gokit/transport/dto"
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
		req := request.(dto.CreateRequest)
		id, err := s.Create(ctx, req.Order)

		return dto.CreateResponse{
			ID:  id,
			Err: err,
		}, nil
	}
}

func makeGetByIDEndpoint(s service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dto.GetByIDRequest)
		orderRequest, err := s.GetByID(ctx, req.ID)

		return dto.GetByIDResponse{
			Order: orderRequest,
			Err:   err,
		}, nil
	}
}

func makeChangeStatusEndpoint(s service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dto.ChangeStatusRequest)
		err = s.ChangeStatus(ctx, req.ID, req.Status)

		return dto.ChangeStatusResponse{Err: err}, nil
	}
}
