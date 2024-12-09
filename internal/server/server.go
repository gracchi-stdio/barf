package server

import (
	"fmt"
	"github.com/gracchi-stdio/barf/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	e   *echo.Echo
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Str("user-agent", c.Request().UserAgent()).
				Msg("request")
			return nil
		},
	}))
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	return &Server{
		e:   e,
		cfg: cfg,
	}
}

func (s *Server) Run() {
	log.Info().
		Str("port", s.cfg.Server.Port).
		Msg("starting server")
	if err := s.e.Start(fmt.Sprintf(":%v", s.cfg.Server.Port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
