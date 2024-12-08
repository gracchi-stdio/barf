package server

import (
	"fmt"
	"github.com/gracchi-stdio/barf/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Server struct {
	e      *echo.Echo
	cfg    *config.Config
	logger zerolog.Logger
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

}

func New(cfg *config.Config) *Server {
	logger := log.With().Str("component", "server").Logger()

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
	}))
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	return &Server{
		e:      e,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() {
	s.e.Logger.Fatal(s.e.Start(fmt.Sprintf(":%v", s.cfg.Server.Port)))
}
