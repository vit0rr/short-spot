package users

import (
	"net/http"

	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/log"
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
	users, err := h.service.List(r.Context(), *h.service.deps.DBClient)
	if err != nil {
		log.Error(r.Context(), "Failed to list users", log.ErrAttr(err))
		return nil, err
	}

	return users, nil
}

// POST /users/create
func (h *HTTP) Create(_ http.ResponseWriter, r *http.Request) (interface{}, error) {

	success, err := h.service.Create(r.Context(), r.Body, *h.service.deps.DBClient)
	if err != nil {
		log.Error(r.Context(), "Failed to create user", log.ErrAttr(err))
		return nil, err
	}

	return success, nil
}
