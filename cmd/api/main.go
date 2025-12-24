package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/touros-platform/api/internal/config"
	"github.com/touros-platform/api/internal/database"
	"github.com/touros-platform/api/internal/handler"
	"github.com/touros-platform/api/internal/observability"
	"github.com/touros-platform/api/internal/repository"
	"github.com/touros-platform/api/internal/router"
	"github.com/touros-platform/api/internal/service"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := observability.NewLogger(cfg.App.Environment, cfg.App.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	var tp *tracesdk.TracerProvider
	if cfg.OTEL.Enabled {
		tp, err = observability.InitTracer(cfg.OTEL.ServiceName, cfg.OTEL.Endpoint)
		if err != nil {
			logger.Warn("Failed to initialize tracing", zap.Error(err))
		} else {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := tp.Shutdown(ctx); err != nil {
					logger.Error("Failed to shutdown tracer", zap.Error(err))
				}
			}()
		}
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.AutoMigrate(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	userRepo := repository.NewUserRepository(db)
	agencyRepo := repository.NewAgencyRepository(db)
	guideRepo := repository.NewGuideRepository(db)
	permitRepo := repository.NewPermitRepository(db)
	checkInRepo := repository.NewSafetyCheckInRepository(db)
	incidentRepo := repository.NewIncidentRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	guideService := service.NewGuideService(guideRepo, userRepo)
	agencyService := service.NewAgencyService(agencyRepo)
	permitService := service.NewPermitService(permitRepo, guideRepo)
	safetyService := service.NewSafetyService(checkInRepo, incidentRepo, guideRepo)

	guideHandler := handler.NewGuideHandler(guideService)
	agencyHandler := handler.NewAgencyHandler(agencyService)
	permitHandler := handler.NewPermitHandler(permitService)
	safetyHandler := handler.NewSafetyHandler(safetyService)
	healthHandler := handler.NewHealthHandler(db)

	r := router.SetupRouter(
		cfg,
		logger,
		authService,
		guideHandler,
		agencyHandler,
		permitHandler,
		safetyHandler,
		healthHandler,
	)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

