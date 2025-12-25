package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	loggerPkg "github.com/Rugved7/go-boilerplate/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/Rugved7/go-boilerplate/internal/config"
	"github.com/Rugved7/go-boilerplate/internal/database"
	"github.com/Rugved7/go-boilerplate/internal/lib/job"
)

type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *loggerPkg.LoggerService
	DB            *database.Database
	Redis         *redis.Client
	httpServer    *http.Server
	Job           *job.JobService
}

func New(
	cfg *config.Config,
	logger *zerolog.Logger,
	loggerService *loggerPkg.LoggerService,
) (*Server, error) {

	// Initialize database
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis (NO New Relic hooks — removed upstream)
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	// Test Redis connection (non-fatal)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn().
			Err(err).
			Msg("failed to connect to Redis, continuing without Redis")
	}

	// Initialize job service
	jobService := job.NewJobService(logger, cfg)
	jobService.InitHandlers(cfg, logger)

	if err := jobService.Start(); err != nil {
		return nil, err
	}

	server := &Server{
		Config:        cfg,
		Logger:        logger,
		LoggerService: loggerService,
		DB:            db,
		Redis:         redisClient,
		Job:           jobService,
	}

	return server, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Msg("starting HTTP server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
	}

	if s.DB != nil {
		if err := s.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}

	if s.Redis != nil {
		if err := s.Redis.Close(); err != nil {
			s.Logger.Warn().Err(err).Msg("failed to close Redis connection")
		}
	}

	if s.Job != nil {
		s.Job.Stop()
	}

	return nil
}
