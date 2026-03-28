package router

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"telemetry_bridge/internal/config"
	"telemetry_bridge/internal/router/handler"
	"telemetry_bridge/internal/router/middleware"
	"time"

	log "github.com/dredfort42/go_logger"
	"github.com/gin-gonic/gin"
)

const (
	ShutdownDeadline = 30 * time.Second // Default shutdown deadline for the server to gracefully close connections
)

var server *http.Server

func Init(cancel context.CancelFunc) error {
	setGinMode()

	router := gin.New()

	// Add global middleware in order of execution
	router.Use(middleware.RequestID()) // 1. Correlation ID (must be first so logs/other middleware get the ID)
	router.Use(middleware.SessionID()) // 2. Ensure session id
	router.Use(middleware.Metrics())   // 3. Metrics (wraps the whole request)
	router.Use(middleware.Security())  // 4. Security (handle OPTIONS/preflight before auth)

	router.Use(gin.Logger())   // 5. Access log (can include request ID)
	router.Use(gin.Recovery()) // 6. Panic recovery (after logger so panics get logged)

	publicGroup := router.Group("/")
	publicGroup.Use(middleware.NoCache())
	{
		publicGroup.GET("/health", handler.Health())
		publicGroup.GET("/metrics", handler.Metrics())

		// Ping-pong
		publicGroup.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	privateGroup := router.Group("/")
	// privateGroup.Use(middleware.Auth())
	privateGroup.Use(middleware.Audit())
	{
		// Ping-pong
		privateGroup.GET("/p-ping", func(c *gin.Context) {
			c.String(http.StatusOK, "p-pong")
		})
	}

	// Configure server
	server = &http.Server{
		Addr:           config.App.Service.Host + ":" + strconv.Itoa(config.App.Service.Port),
		Handler:        router,
		ReadTimeout:    config.App.Service.ReadTimeout,
		WriteTimeout:   config.App.Service.WriteTimeout,
		IdleTimeout:    config.App.Service.IdleTimeout,
		MaxHeaderBytes: config.App.Service.MaxHeaderBytes,
	}

	// Start server in a goroutine
	go func(cancel context.CancelFunc) {
		log.Info.Printf("Server starting on %s:%d", config.App.Service.Host, config.App.Service.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error.Printf("Server failed to start: %v", err)
			cancel()
		}
	}(cancel)

	return nil
}

func Close() {
	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownDeadline)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Error.Fatalf("API server forced to shutdown: %v", err)
	}

	log.Info.Println("API server stopped")
}

// setGinMode sets the Gin mode based on the DEBUG environment variable
func setGinMode() {
	if config.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
