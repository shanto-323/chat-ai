package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/database"
	"github.com/shanto-323/chat-ai/internal/server"
	"github.com/shanto-323/chat-ai/internal/server/handler"
	"github.com/shanto-323/chat-ai/internal/server/router"
	"github.com/shanto-323/chat-ai/internal/service"
	logs "github.com/shanto-323/chat-ai/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("error loading config " + err.Error())
	}

	logger := logs.New(cfg)

	if err := database.Migrate(context.Background(), logger, cfg); err != nil {
		logger.Fatal().Err(err).Msg("failed to migrate database")
	}

	server, err := server.New(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize server")
	}

	service := service.New(server)
	handler := handler.New(server, service)

	router := router.NewRouter(server, handler)
	server.SetUpHTTPServer(router)

	stopChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := server.Run(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-stopChan:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- server.Stop(ctx)
		}()

		select {
		case <-ctx.Done():
			logger.Fatal().Err(ctx.Err()).Msg("timeout stopping server")
		case err := <-done:
			if err != nil {
				logger.Fatal().Err(err).Msg("error stopping server")
			}
			logger.Info().Msg("server stopped")
		}
	case err := <-errChan:
		logger.Fatal().Err(err).Msg("error running server " + err.Error())
	}
}
