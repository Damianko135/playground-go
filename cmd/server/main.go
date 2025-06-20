package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Damianko135/playground-go/internal/config"
	"github.com/Damianko135/playground-go/internal/handlers"
	"github.com/Damianko135/playground-go/internal/middleware"
	"github.com/Damianko135/playground-go/internal/utils"
	"github.com/Damianko135/playground-go/views"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("❌ Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Print configuration
	cfg.Print()

	fmt.Println("🔧 Starting Echo server...")
	e := echo.New()

	// Hide Echo banner
	e.HideBanner = true
	e.Debug = cfg.Server.Debug

	// Apply core middleware
	e.Use(echomiddleware.Recover())
	e.Use(middleware.SecurityHeaders())
	e.Use(middleware.RequestID())
	e.Use(middleware.ResponseTime())

	// Conditional middleware based on configuration
	if cfg.IsDevelopment() {
		e.Use(middleware.CustomLogger())
		fmt.Println("🐛 Debug mode enabled")
	}

	if cfg.API.EnableCORS {
		e.Use(middleware.CORS())
	}

	if cfg.API.EnableGzip {
		e.Use(middleware.Gzip())
	}

	e.Use(middleware.Cache())

	// Metrics middleware (always enabled for monitoring)
	e.Use(handlers.MetricsMiddleware())

	// Rate limiting for API endpoints
	apiGroup := e.Group("/api")
	apiGroup.Use(middleware.RateLimiter())
	apiGroup.Use(middleware.APIKeyAuth(cfg.API.Key))

	// Static files
	e.Static("/static", "../../static")

	// Web routes
	e.GET("/", utils.Temple(views.Home()))
	e.GET("/about", utils.Temple(views.About()))
	e.GET("/playground", utils.Temple(views.Playground()))
	e.GET("/tools", utils.Temple(views.Tools()))

	// Health check endpoints (if enabled)
	if cfg.Features.EnableHealthCheck {
		e.GET("/health", handlers.HealthCheck)
		e.GET("/health/ready", handlers.ReadinessCheck)
		e.GET("/health/live", handlers.LivenessCheck)
	}

	// Metrics endpoint (if enabled)
	if cfg.Features.EnableMetrics {
		e.GET("/metrics", handlers.GetMetrics)
	}

	// API endpoints
	apiGroup.GET("/weather", handlers.GetWeather)
	apiGroup.GET("/quote", handlers.GetQuote)
	apiGroup.GET("/stats", handlers.GetSystemStats)
	apiGroup.GET("/palette", handlers.GetColorPalette)
	apiGroup.GET("/joke", handlers.GetJoke)
	apiGroup.GET("/random", handlers.GetRandomNumber)
	apiGroup.GET("/timezones", handlers.GetTimeZones)

	// Configure server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Print startup information
	fmt.Printf("🚀 Server starting on port %s\n", cfg.Server.Port)
	if cfg.Features.EnableHealthCheck {
		fmt.Printf("📊 Health check available at: http://localhost:%s/health\n", cfg.Server.Port)
	}
	if cfg.Features.EnableMetrics {
		fmt.Printf("📈 Metrics available at: http://localhost:%s/metrics\n", cfg.Server.Port)
	}
	fmt.Printf("🎮 Playground available at: http://localhost:%s/playground\n", cfg.Server.Port)
	fmt.Printf("🔧 Tools available at: http://localhost:%s/tools\n", cfg.Server.Port)
	fmt.Printf("📡 API endpoints available at: http://localhost:%s/api/*\n", cfg.Server.Port)

	// Start server in a goroutine
	go func() {
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("✅ Server gracefully stopped")
}
