package main

import (
	"github.com/gracchi-stdio/barf/internal/config"
	"github.com/gracchi-stdio/barf/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	server := server.New(cfg)

	server.Run()
}
