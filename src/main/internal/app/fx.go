package app

import (
	"context"
	"os"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/apikey"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/middleware"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// NewConfig loads and returns application config using Viper
func NewConfig() *config.Config {
	return config.Load()
}

// NewGinEngine creates and configures the Gin engine
func NewGinEngine(cfg *config.Config) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CorrelationMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	r.NoMethod(func(c *gin.Context) {
		response.New(c).
			Error(response.MethodNotAllowed, "Method not allowed: "+c.Request.Method).
			WithMeta().
			JSONAbort()
	})

	r.NoRoute(func(c *gin.Context) {
		response.New(c).
			Error(response.NotFound, "Route not found: "+c.Request.URL.Path).
			WithMeta().
			JSONAbort()
	})

	r.HandleMethodNotAllowed = true
	return r
}

// NewRouterGroup creates the main router group with a context path
func NewRouterGroup(r *gin.Engine, cfg *config.Config, db *gorm.DB) *types.RouterGroup {
	log.Info().Str("contextPath", cfg.App.ContextPath).Msg("Server context path")

	// Set DB for apikey middleware
	apikey.SetDB(db)

	def := r.Group(cfg.App.ContextPath)

	// Health check endpoint (simple ping)
	def.GET("/ping", func(c *gin.Context) {
		response.OK(c, "Success ping", cfg.App.Name)
	})

	// Database health check endpoint with correlation ID in response
	def.GET("/health", func(c *gin.Context) {
		if err := database.HealthCheck(db); err != nil {
			response.InternalError(c, "Database unhealthy", err.Error())
			return
		}

		stats, _ := database.Stats(db)
		response.OK(c, "Service healthy", gin.H{
			"database": gin.H{
				"status":       "connected",
				"openConns":    stats.OpenConnections,
				"inUse":        stats.InUse,
				"idle":         stats.Idle,
				"maxOpen":      stats.MaxOpenConnections,
				"waitCount":    stats.WaitCount,
				"waitDuration": stats.WaitDuration.String(),
			},
		})
	})

	return &types.RouterGroup{RouterGroup: def}
}

// ServerParams contains dependencies for starting the server
type ServerParams struct {
	fx.In

	Engine      *gin.Engine
	Config      *config.Config
	RouterGroup *types.RouterGroup
	Lifecycle   fx.Lifecycle
}

// RegisterServer registers the HTTP server with fx lifecycle
func RegisterServer(p ServerParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info().Str("port", p.Config.App.Port).Msg("🚀 Server running")
				if err := p.Engine.Run(":" + p.Config.App.Port); err != nil {
					log.Error().Err(err).Msg("Server failed to start")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Shutting down server...")
			return nil
		},
	})
}

// initializeApp performs pre-fx initialization
func initializeApp() {
	// Load config first
	cfg := config.Load()

	// Initialize logger
	if cfg.App.Env != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.App.Env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Str("env", cfg.App.Env).Msg("Application initialized")

	// Setup timezone
	log.Info().Str("timezone", cfg.App.Timezone).Msg("Loading timezone")
	loc, err := time.LoadLocation(cfg.App.Timezone)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load timezone")
	}
	time.Local = loc

	// Register event handlers for cache invalidation
	observer.RegisterEvent()
}

// CoreModule provides core infrastructure
var CoreModule = fx.Module("core",
	fx.Provide(
		NewConfig,
		database.NewDatabase, // Use FX-aware database provider
		NewGinEngine,
		NewRouterGroup,
	),
	fx.Invoke(RegisterServer),
)
