package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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
	r.Use(SetResponseTypeToJSON)
	r.Use(telemetry.TelemetryMiddleware)

	// Custom middleware
	r.Use(corsMiddleware)
	r.Use(authMiddleware)

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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	godotenv.Load()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != os.Getenv("AUTH_TOKEN") {
			http.Error(w, "Unauthorized. Please provide a valid auth token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
