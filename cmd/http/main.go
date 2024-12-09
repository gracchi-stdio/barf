package main

import (
	"github.com/gracchi-stdio/barf/internal/config"
	"github.com/gracchi-stdio/barf/internal/server"
	"github.com/gracchi-stdio/barf/pkg/logger"
)

func main() {

	logger.Setup(logger.Config{
		Environment: "production",
		LogLevel:    "debug",
	})
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	server := server.New(cfg)

	server.Run()
}
