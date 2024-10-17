package server

import (
	"context"
	"net/http"
	"time"

	"github.com/vit0rr/short-spot/api/router"
	"github.com/vit0rr/short-spot/pkg/deps"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(ctx context.Context, deps *deps.Deps, db *mongo.Database) *http.Server {
	router := router.New(deps, *db)
	timeout := time.Duration(deps.Config.Server.CtxTimeout) * time.Second

	return &http.Server{
		Addr:         deps.Config.Server.BindAddr,
		Handler:      http.TimeoutHandler(router.BuildRoutes(), timeout, ""),
		ReadTimeout:  timeout,
		WriteTimeout: timeout + 2,
	}
}
