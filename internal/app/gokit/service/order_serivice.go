package service

import (
	"context"
	"github.com/AlexBBBril/gokit/internal/app/gokit/entity"
)

// OrderService describes the Order service.
type OrderService interface {
	Create(ctx context.Context, order entity.Order) (string, error)
	GetByID(ctx context.Context, id string) (entity.Order, error)
	ChangeStatus(ctx context.Context, id, status string) error
}
