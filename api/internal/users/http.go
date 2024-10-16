package users

import (
	"net/http"

	"github.com/vit0rr/short-spot/pkg/deps"
)

type HTTP struct {
	service *Service
}

func NewHTTP(deps *deps.Deps) *HTTP {
	return &HTTP{
		service: NewService(deps),
	}
}

// GET /users
func (h *HTTP) List(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	return h.service.List(r.Context())
}
