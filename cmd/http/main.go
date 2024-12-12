package main

import (
	"github.com/gracchi-stdio/barf/internal/config"
	"github.com/gracchi-stdio/barf/internal/server"
	"github.com/gracchi-stdio/barf/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err)
	}

	logger.Setup(logger.Config{
		Environment: logger.Environment(cfg.Server.Env),
		LogLevel:    "debug",
	})

	server := server.New(cfg)

	if err := server.Run(); err != nil {
		log.Fatal().Err(err)
	}
}
