package router

import (
	"github.com/gin-gonic/gin"
	"github.com/touros-platform/api/internal/config"
	"github.com/touros-platform/api/internal/handler"
	"github.com/touros-platform/api/internal/middleware"
	"github.com/touros-platform/api/internal/service"
	"go.uber.org/zap"
)

func SetupRouter(
	cfg *config.Config,
	logger *zap.Logger,
	authService service.AuthService,
	guideHandler *handler.GuideHandler,
	agencyHandler *handler.AgencyHandler,
	permitHandler *handler.PermitHandler,
	safetyHandler *handler.SafetyHandler,
	healthHandler *handler.HealthHandler,
) *gin.Engine {
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.RateLimitMiddleware(100, 200))

	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
	SetupMetrics(r)

	auth := r.Group("/api/v1/auth")
	{
		authHandler := handler.NewAuthHandler(authService)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(authService))
	{
		guides := api.Group("/guides")
		{
			guides.POST("", guideHandler.Create)
			guides.GET("", guideHandler.List)
			guides.GET("/:id", guideHandler.GetByID)
			guides.PUT("/:id", guideHandler.Update)
			guides.POST("/:id/verify", middleware.RequireRole("admin"), guideHandler.Verify)
			guides.POST("/:id/suspend", middleware.RequireRole("admin"), guideHandler.Suspend)
		}

		agencies := api.Group("/agencies")
		{
			agencies.POST("", agencyHandler.Create)
			agencies.GET("", agencyHandler.List)
			agencies.GET("/:id", agencyHandler.GetByID)
			agencies.PUT("/:id", agencyHandler.Update)
			agencies.POST("/:id/verify", middleware.RequireRole("admin"), agencyHandler.Verify)
			agencies.POST("/:id/suspend", middleware.RequireRole("admin"), agencyHandler.Suspend)
		}

		permits := api.Group("/permits")
		{
			permits.POST("", permitHandler.Create)
			permits.GET("", permitHandler.List)
			permits.GET("/:id", permitHandler.GetByID)
			permits.POST("/:id/revoke", middleware.RequireRole("admin"), permitHandler.Revoke)
		}

		permitsPublic := r.Group("/api/v1/permits")
		{
			permitsPublic.GET("/validate/:number", permitHandler.Validate)
		}

		safety := api.Group("/safety")
		{
			safety.POST("/check-ins", safetyHandler.CreateCheckIn)
			safety.GET("/check-ins/:id", safetyHandler.GetCheckInByID)
			safety.GET("/guides/:guide_id/check-ins", safetyHandler.ListCheckIns)

			safety.POST("/incidents", safetyHandler.CreateIncident)
			safety.GET("/incidents", safetyHandler.ListIncidents)
			safety.GET("/incidents/:id", safetyHandler.GetIncidentByID)
			safety.PUT("/incidents/:id", safetyHandler.UpdateIncident)
			safety.GET("/guides/:guide_id/sos", safetyHandler.GetActiveSOS)
		}
	}

	return r
}

