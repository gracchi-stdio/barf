package server

import (
	"fmt"
	"github.com/gracchi-stdio/barf/internal/config"
	"github.com/gracchi-stdio/barf/internal/domain"
	httphandler "github.com/gracchi-stdio/barf/internal/handler/http"
	"github.com/gracchi-stdio/barf/internal/repository"
	"github.com/gracchi-stdio/barf/internal/service"
	"github.com/gracchi-stdio/barf/pkg/bookfetcher"
	"github.com/gracchi-stdio/barf/pkg/bookfetcher/providers/googlebooks"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Server struct {
	e   *echo.Echo
	cfg *config.Config
	db  *gorm.DB
}

func New(cfg *config.Config) *Server {
	e := echo.New()

	// Request logging middleware
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
	e.Use(middleware.CORS())

	return &Server{
		e:   e,
		cfg: cfg,
	}
}

func (s *Server) setupDB() error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		s.cfg.DB.Host,
		s.cfg.DB.User,
		s.cfg.DB.Password,
		s.cfg.DB.Name,
		s.cfg.DB.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// enable uuid
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return fmt.Errorf("failed to create uuid extension: %w", err)
	}

	// migrate database
	if err := db.AutoMigrate(
		domain.Book{},
		domain.Inventory{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Info().Msg("database migrated")

	s.db = db

	return nil
}
func (s *Server) setupRoutes() {

	// initialize repositories
	bookRepo := repository.NewBookRepository(s.db)
	inventoryRepo := repository.NewInventoryRepository(s.db)

	// google books fetcher
	googleProvider := googlebooks.NewGoogleBooksProvider("", 5*time.Second)

	// initialize fetchers
	fetchers := map[string]bookfetcher.BookFetcher{
		"googlebooks": googleProvider,
	}
	// initialize services
	bookService := service.NewBookService(
		bookRepo,
		inventoryRepo,
		fetchers,
		"googlebooks")

	// initialize handlers
	bookHandler := httphandler.NewBookHandler(bookService)

	bookHandler.RegisterRoutes(s.e)

	s.e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
}

func (s *Server) Run() error {
	if err := s.setupDB(); err != nil {
		return err
	}

	s.setupRoutes()

	log.Info().
		Str("port", s.cfg.Server.Port).
		Msg("starting server")

	if err := s.e.Start(fmt.Sprintf(":%v", s.cfg.Server.Port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}

	return nil
}
