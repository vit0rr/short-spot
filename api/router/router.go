package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vit0rr/short-spot/api/internal/users"
	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/telemetry"
)

type Router struct {
	Deps  *deps.Deps
	users *users.HTTP
}

func New(deps *deps.Deps) *Router {
	return &Router{
		Deps:  deps,
		users: users.NewHTTP(deps),
	}
}

func (router *Router) BuildRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)
	r.Use(SetResponseTypeToJSON)
	r.Use(telemetry.TelemetryMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", telemetry.HandleFuncLogger(router.users.List))
			r.Post("/create", telemetry.HandleFuncLogger(router.users.Create))
		})
	})

	return r
}
