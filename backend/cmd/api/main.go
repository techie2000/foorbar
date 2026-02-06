package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/techie2000/axiom/internal/config"
	"github.com/techie2000/axiom/internal/handler"
	"github.com/techie2000/axiom/internal/middleware"
	"github.com/techie2000/axiom/internal/repository"
	"github.com/techie2000/axiom/internal/service"
	"github.com/techie2000/axiom/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/techie2000/axiom/docs" // Swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Axiom API
// @version 1.0
// @description Financial Services Static Data Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@axiom.example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.Log.Level)

	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos)

	// Initialize handlers
	handlers := handler.NewHandlers(services)

	// Setup Gin router
	router := setupRouter(cfg, handlers)

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info().Msgf("Starting Axiom API server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info().Msg("Server exited")
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info().Msg("Database connection established")
	return db, nil
}

func setupRouter(cfg *config.Config, h *handler.Handlers) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.RateLimit())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Auth.Login)
			auth.POST("/register", h.Auth.Register)
		}

		// Protected routes (require JWT)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg))
		{
			// Domain data routes
			countries := protected.Group("/countries")
			{
				countries.GET("", h.Country.List)
				countries.GET("/:id", h.Country.Get)
				countries.POST("", h.Country.Create)
				countries.PUT("/:id", h.Country.Update)
				countries.DELETE("/:id", h.Country.Delete)
			}

			currencies := protected.Group("/currencies")
			{
				currencies.GET("", h.Currency.List)
				currencies.GET("/:id", h.Currency.Get)
				currencies.POST("", h.Currency.Create)
				currencies.PUT("/:id", h.Currency.Update)
				currencies.DELETE("/:id", h.Currency.Delete)
			}

			entities := protected.Group("/entities")
			{
				entities.GET("", h.Entity.List)
				entities.GET("/:id", h.Entity.Get)
				entities.POST("", h.Entity.Create)
				entities.PUT("/:id", h.Entity.Update)
				entities.DELETE("/:id", h.Entity.Delete)
			}

			instruments := protected.Group("/instruments")
			{
				instruments.GET("", h.Instrument.List)
				instruments.GET("/:id", h.Instrument.Get)
				instruments.POST("", h.Instrument.Create)
				instruments.PUT("/:id", h.Instrument.Update)
				instruments.DELETE("/:id", h.Instrument.Delete)
			}

			accounts := protected.Group("/accounts")
			{
				accounts.GET("", h.Account.List)
				accounts.GET("/:id", h.Account.Get)
				accounts.POST("", h.Account.Create)
				accounts.PUT("/:id", h.Account.Update)
				accounts.DELETE("/:id", h.Account.Delete)
			}

			ssis := protected.Group("/ssis")
			{
				ssis.GET("", h.SSI.List)
				ssis.GET("/:id", h.SSI.Get)
				ssis.POST("", h.SSI.Create)
				ssis.PUT("/:id", h.SSI.Update)
				ssis.DELETE("/:id", h.SSI.Delete)
			}

			// Data acquisition routes
			dataAcq := protected.Group("/data")
			{
				dataAcq.POST("/import", h.DataAcquisition.Import)
				dataAcq.POST("/export", h.DataAcquisition.Export)
				dataAcq.GET("/jobs", h.DataAcquisition.ListJobs)
				dataAcq.GET("/jobs/:id", h.DataAcquisition.GetJob)
			}
		}
	}

	return router
}
