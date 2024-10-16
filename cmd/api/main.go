package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vit0rr/short-spot/api/server"
	"github.com/vit0rr/short-spot/config"
	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/log"
)

func main() {
	ctx := context.Background()
	// Parse config file location from command-line flag
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the config file (see config/config_local.hcl for an example)")
	flag.Parse()

	var cfg config.Config
	if configPath == "" {
		cfg = config.DefaultConfig()
	} else if configPath != "" {
		parsedConfig, err := config.GetConfig(configPath)
		if err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			os.Exit(1)
		}
		cfg = parsedConfig
	}

	logLevel, err := log.ParseLogLevel(cfg.Server.LogLevel)
	if err != nil {
		fmt.Printf("Error parsing log level: %v\n`Info` is applied as a default\n", err)
		logLevel = slog.LevelInfo
	}
	log.New(ctx, logLevel)

	pgCfg, err := pgxpool.ParseConfig(cfg.API.Postgres.Dsn)
	if err != nil {
		log.Error(ctx, "unable to parse database connection", log.ErrAttr(err))
		os.Exit(1)
	}

	// initialize postgres pool
	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		log.Error(ctx, "unable to initialize database connection", log.ErrAttr(err))
		os.Exit(1)
	}
	defer pool.Close()

	dependencies := deps.New(cfg, pool)

	httpServer := server.New(ctx, dependencies)

	// Handle graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error(ctx, "unexpected error during server shutdown", log.ErrAttr(err))
		}
		close(idleConnsClosed)
	}()

	log.Info(ctx, "Starting API at", log.AnyAttr("bind_addr", cfg.Server.BindAddr))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(ctx, "error starting server", log.ErrAttr(err))
	}

	<-idleConnsClosed
}
