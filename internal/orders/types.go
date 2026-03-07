package orders

import (
	"context"

	repo "github.com/juxtapsy2/go_ecom/internal/adapters/postgresql/sqlc"
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}

type orderItem struct {
	ProductID    int64 `json:"productId"`
	Quantity     int32 `json:"quantity"`
	PriceInCents int32 `json:"price_cents"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}
