package dto

import "github.com/AlexBBBril/gokit/internal/app/gokit/entity"

type (
	// CreateRequest holds the request parameters for the Create method.
	CreateRequest struct {
		Order entity.Order
	}

	// CreateResponse holds the response values for the Create method.
	CreateResponse struct {
		ID  string `json:"id"`
		Err error  `json:"error, omitempty"`
	}

	// GetByIDRequest holds the request parameters for the GetByID method.
	GetByIDRequest struct {
		ID string
	}

	// GetByIDResponse holds the response values for the GetByID method.
	GetByIDResponse struct {
		Order entity.Order `json:"order"`
		Err   error        `json:"error, omitempty"`
	}

	// ChangeStatusRequest holds the request parameters for the ChangeStatus method.
	ChangeStatusRequest struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	// ChangeStatusResponse holds the response values for the ChangeStatus method.
	ChangeStatusResponse struct {
		Err error `json:"error,omitempty"`
	}
)
