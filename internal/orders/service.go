package orders

import (
	"context"
	"errors"
	"time"

	repo "github.com/AmmanSajid1/go-ecom-api/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidOrder    = errors.New("order must contain at least one item")
	ErrInvalidCustomer = errors.New("customer_id must be greater than 0")
	ErrInvalidItem     = errors.New("each item must have product_id, quantity and price_cents greater than 0")
	ErrProductNoStock  = errors.New("product does not have enough stock")
)

type Service interface {
	PlaceOrder(ctx context.Context, customerId int, items []OrderItemInput) (OrderResponse, error)
	GetOrderByID(ctx context.Context, orderId int) (OrderResponse, error)
}

type svc struct {
	db   *pgx.Conn
	repo *repo.Queries
}

func NewService(db *pgx.Conn) Service {
	return &svc{
		db:   db,
		repo: repo.New(db),
	}
}

func (s *svc) PlaceOrder(
	ctx context.Context,
	customerId int,
	items []OrderItemInput,
) (OrderResponse, error) {
	if customerId <= 0 {
		return OrderResponse{}, ErrInvalidCustomer
	}
	if len(items) == 0 {
		return OrderResponse{}, ErrInvalidOrder
	}
	for _, item := range items {
		if item.ProductID <= 0 || item.Quantity <= 0 || item.PriceCents <= 0 {
			return OrderResponse{}, ErrInvalidItem
		}
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return OrderResponse{}, err
	}

	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, int64(customerId))
	if err != nil {
		return OrderResponse{}, err
	}

	for _, item := range items {
		rowsAffected, err := qtx.DecrementProductStock(ctx, repo.DecrementProductStockParams{
			Quantity: int32(item.Quantity),
			ID:       int64(item.ProductID),
		})
		if err != nil {
			return OrderResponse{}, err
		}

		if rowsAffected == 0 {
			return OrderResponse{}, ErrProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  int64(item.ProductID),
			Quantity:   int32(item.Quantity),
			PriceCents: int32(item.PriceCents),
		})
		if err != nil {
			return OrderResponse{}, err
		}
	}

	orderItems, err := qtx.ListOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return OrderResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return OrderResponse{}, err
	}

	responseItems := make([]OrderItemResponse, 0, len(orderItems))

	for _, item := range orderItems {
		responseItems = append(responseItems, OrderItemResponse{
			ProductID:  int(item.ProductID),
			Quantity:   int(item.Quantity),
			PriceCents: int(item.PriceCents),
		})
	}

	return OrderResponse{
		ID:         int(order.ID),
		CustomerID: int(order.CustomerID),
		CreatedAt:  order.CreatedAt.Time.Format(time.RFC3339),
		Items:      responseItems,
	}, nil
}

func (s *svc) GetOrderByID(ctx context.Context, orderId int) (OrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, int64(orderId))
	if err != nil {
		return OrderResponse{}, err
	}

	orderItems, err := s.repo.ListOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return OrderResponse{}, err
	}

	responseItems := make([]OrderItemResponse, 0, len(orderItems))

	for _, item := range orderItems {
		responseItems = append(responseItems, OrderItemResponse{
			ProductID:  int(item.ProductID),
			Quantity:   int(item.Quantity),
			PriceCents: int(item.PriceCents),
		})
	}

	return OrderResponse{
		ID:         int(order.ID),
		CustomerID: int(order.CustomerID),
		CreatedAt:  order.CreatedAt.Time.Format(time.RFC3339),
		Items:      responseItems,
	}, nil
}
