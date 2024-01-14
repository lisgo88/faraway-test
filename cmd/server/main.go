package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/lisgo88/faraway-test/internal/app"
	"github.com/lisgo88/faraway-test/internal/config"
	"github.com/lisgo88/faraway-test/internal/pkg/repository/quotes"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init logger
	logger := zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		},
	).With().Timestamp().Logger()

	// get config
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("can't get config")
	}

	logger.Info().Any("msg", cfg).Msg("server config")
	logger = logger.Level(cfg.LogLevel) // set log level

	// quotes repository
	quotesRepo := quotes.New(ctx, logger)

	// init service
	server := app.NewServer(quotesRepo, cfg.Server, logger)

	// run tcp server
	err = server.Run(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("server error")
	}
}
