package cockroachdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AlexBBBril/gokit/internal/app/gokit/entity"
	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var ErrRepository = errors.New("unable to handle request")

type repository struct {
	db     *sql.DB
	logger log.Logger
}

// New returns a concrete repository backed by CockroachDB
func New(db *sql.DB, logger log.Logger) (entity.OrderRepository, error) {
	// return  repository
	return &repository{
		db:     db,
		logger: log.With(logger, "rep", "cockroachdb"),
	}, nil
}

// CreateOrder inserts a new order and its order items into db
func (r repository) CreateOrder(ctx context.Context, order entity.Order) error {
	// Run a transaction to sync the query model.
	err := crdb.ExecuteTx(ctx, r.db, nil, func(tx *sql.Tx) error {
		return r.createOrder(tx, order)
	})
	if err != nil {
		return err
	}
	return nil
}

// ChangeOrderStatus changes the order status
func (r repository) createOrder(tx *sql.Tx, order entity.Order) error {
	query := `INSERT INTO orders (id, customerid, status, createdon, restaurantid)
			VALUES ($1,$2,$3,$4,$5)`

	_, err := tx.Exec(query, order.ID, order.CustomerID, order.Status, order.CreatedOn, order.RestaurantID)
	if nil != err {
		return err
	}

	for _, v := range order.OrderItems {
		query = `INSERT INTO orderitems (orderid, customerid, code, name, unitprice, quantity)
			   VALUES ($1,$2,$3,$4,$5,$6)`

		_, err := tx.Exec(query, order.ID, order.CustomerID, v.ProductCode, v.Name, v.UnitPrice, v.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r repository) ChangeOrderStatus(ctx context.Context, id string, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE orders SET status=$2 WHERE id=$1`, id, status)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderByID query the order by given id
func (r repository) GetOrderByID(ctx context.Context, id string) (entity.Order, error) {
	var orderRow = entity.Order{}
	query := "SELECT id, customerid, status, createdon, restaurantid FROM orders WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&orderRow.ID, &orderRow.CustomerID, &orderRow.Status, &orderRow.CreatedOn, &orderRow.RestaurantID)

	if nil != err {
		_ = level.Error(r.logger).Log("err", err.Error())
		return orderRow, err
	}

	return orderRow, nil
}

// Close implements DB.Close
func (r *repository) Close() error {
	return r.db.Close()
}
