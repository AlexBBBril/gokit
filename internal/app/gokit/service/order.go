package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AlexBBBril/gokit/internal/app/gokit/entity"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
	"time"
)

// Service error
var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)

type (
	// OrderService describes the Order service.
	OrderService interface {
		Create(ctx context.Context, order entity.Order) (string, error)
		GetByID(ctx context.Context, id string) (entity.Order, error)
		ChangeStatus(ctx context.Context, id, status string) error
	}

	// orderService implements the Order Service
	orderService struct {
		repository entity.OrderRepository
		logger     log.Logger
	}
)

// NewService creates and returns a new Order service instance
func NewService(repository entity.OrderRepository, logger log.Logger) OrderService {
	return &orderService{
		repository: repository,
		logger:     logger,
	}
}

// Create makes an order
func (s *orderService) Create(ctx context.Context, order entity.Order) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	uuid, err := uuid.NewV4()
	if nil != err {
		_ = level.Error(logger).Log("err", err)

		return "", err
	}

	id := uuid.String()
	order.ID = id
	order.Status = "Pending"
	order.CreatedOn = time.Now().Unix()

	if err := s.repository.CreateOrder(ctx, order); nil != err {
		_ = level.Error(logger).Log("err", err)

		return "", err
	}

	return id, nil
}

// GetByID returns an order given by id
func (s *orderService) GetByID(ctx context.Context, id string) (entity.Order, error) {
	logger := log.With(s.logger, "method", "GetByID")
	order, err := s.repository.GetOrderByID(ctx, id)

	if nil != err {
		_ = level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return order, ErrOrderNotFound
		}

		return order, ErrQueryRepository
	}

	return order, nil
}

// ChangeStatus changes the status of an order
func (s *orderService) ChangeStatus(ctx context.Context, id, status string) error {
	logger := log.With(s.logger, "method", "ChangeStatus")
	if err := s.repository.ChangeOrderStatus(ctx, id, status); nil != err {
		_ = level.Error(logger).Log("err", err)
		return ErrCmdRepository
	}

	return nil
}
