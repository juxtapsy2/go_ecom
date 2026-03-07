package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/juxtapsy2/go_ecom/internal/adapters/postgresql/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found.")
	ErrProductNoStock  = errors.New("not enough stock")
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	// validate payload
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required.")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("there must be at least one item in cart to place order.")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	// look for the product if exists
	for _, item := range tempOrder.Items {
		product, err := s.repo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}
		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: item.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	tx.Commit(ctx)

	return order, nil
}
