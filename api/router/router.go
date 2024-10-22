package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	urlshort "github.com/vit0rr/short-spot/api/internal/url-short"
	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/telemetry"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	Deps     *deps.Deps
	urlshort *urlshort.HTTP
}

func (router *Router) BuildRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)

	// Custom middlewares
	r.Use(telemetry.TelemetryMiddleware)
	r.Use(SetResponseTypeToJSON)
	r.Use(CorsMiddleware)
	r.Use(AuthMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", telemetry.HandleFuncLogger(router.urlshort.Redirect))
		r.Route("/short-url", func(r chi.Router) {
			r.Post("/", telemetry.HandleFuncLogger(router.urlshort.ShortUrl))
		})

	})

	return r
}

func New(deps *deps.Deps, db mongo.Database) *Router {
	return &Router{
		Deps:     deps,
		urlshort: urlshort.NewHTTP(deps, &db),
	}
}
