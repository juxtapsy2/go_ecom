package products

import (
	"context"

	repo "github.com/juxtapsy2/go_ecom/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(r repo.Querier) Service {
	return &svc{repo: r}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}
