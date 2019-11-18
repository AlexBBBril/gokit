package entity

import "context"

type (
	// Order represents an order
	Order struct {
		ID           string      `json:"id,omitempty"`
		CustomerID   string      `json:"customer_id"`
		Status       string      `json:"status"`
		CreatedOn    int64       `json:"created_on,omitempty"`
		RestaurantID string      `json:"restaurant_id"`
		OrderItems   []OrderItem `json:"order_items,omitempty"`
	}

	// OrderItem represents orders item
	OrderItem struct {
		ProductCode string  `json:"product_code"`
		Name        string  `json:"name"`
		UnitPrice   float32 `json:"unit_price"`
		Quantity    int32   `json:"quantity"`
	}

	// OrderRepository describes the persistence on order model
	OrderRepository interface {
		CreateOrder(ctx context.Context, order Order) error
		GetOrderByID(ctx context.Context, id string) (Order, error)
		ChangeOrderStatus(ctx context.Context, id string, status string) error
	}
)
