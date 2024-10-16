package users

import (
	"context"

	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/log"
)

type Service struct {
	deps *deps.Deps
}

func NewService(deps *deps.Deps) *Service {
	return &Service{
		deps: deps,
	}
}

func (s *Service) List(ctx context.Context) ([]map[string]interface{}, error) {
	log.Info(ctx, "Listing users inside the service hurray!!")
	return []map[string]interface{}{
		{"id": "123"},
	}, nil
}
