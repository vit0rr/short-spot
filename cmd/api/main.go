package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/vit0rr/short-spot/api/server"
	_ "github.com/vit0rr/short-spot/cmd/api/docs"
	"github.com/vit0rr/short-spot/config"
	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//	@title			Short Spot API
//	@version		1.0
//	@description	Short Spot API is a simple URL shortener service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
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

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.API.Mongo.Dsn))
	if err != nil {
		log.Error(ctx, "unable to parse database connection", log.ErrAttr(err))
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Error(ctx, "unable to disconnect from MongoDB", log.ErrAttr(err))
			panic(err)
		}
	}()

	dependencies := deps.New(cfg, mongoClient)

	httpServer := server.New(ctx, dependencies, mongoClient.Database("shortspot"))

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
