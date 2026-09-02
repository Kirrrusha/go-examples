package main

import (
	"context"

	"account/internal/app"
	"account/internal/config"
	"account/internal/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	l := logger.New(cfg)
	application := app.New(&l, cfg)
	defer application.Close()

	if err := application.Run(ctx); err != nil {
		l.Fatal().Err(err).Msg("error")
	}
}
