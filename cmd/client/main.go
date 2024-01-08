package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/lisgo88/faraway-test/internal/app"
	"github.com/lisgo88/faraway-test/internal/config"
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

	logger.Info().Any("msg", cfg).Msg("client config")
	logger = logger.Level(cfg.LogLevel) // set log level

	// init client
	client := app.NewClient(cfg.Client, logger)

	// run tcp client
	err = client.Run(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("client error")
	}
}
