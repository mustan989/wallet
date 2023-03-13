package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	. "github.com/mustan989/wallet/app/config"
	"github.com/mustan989/wallet/pkg/config"
	"github.com/mustan989/wallet/pkg/logger"
	"github.com/mustan989/wallet/pkg/postgres"
	"github.com/mustan989/wallet/pkg/shutdown"
)

// TODO: mv to env
const configPath = "configs/sandbox.yaml"

func main() {
	ctx := context.Background()

	log := logger.NewLogger(
		logger.WithOutputFormat(logger.JSON),
		logger.WithSeverityWriter(logger.Error, os.Stderr),
	)

	log.Infof("Starting app")
	log.Infof("Getting config")

	var cfg Config

	if err := config.ParseConfigFile(configPath, &cfg); err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	log.Infof("Config successfully loaded")

	connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Infof("Connecting to database")

	pool, err := postgres.Connect(connCtx, cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Error pinging database: %s", err)
	}

	log.Infof("Successfully connected to database")

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	log.Infof("Starting server on port :%d", cfg.Server.Port)

	go func() {
		<-shutdown.GracefulShutdown(
			ctx,
			map[string]shutdown.Operation{
				"Database": func(_ context.Context) error {
					pool.Close()
					return nil
				},
				"Server": func(ctx context.Context) error {
					return e.Shutdown(ctx)
				},
			},
			shutdown.WithLogger(log), shutdown.WithTimeout(1*time.Second),
		)
	}()

	if err = e.Start(fmt.Sprint(":", cfg.Server.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("Error running server: %s", err)
	}

	log.Infof("App finished successfully")
}
