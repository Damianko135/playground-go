package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	System    SystemInfo        `json:"system"`
	Checks    map[string]string `json:"checks"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCPU       int    `json:"num_cpu"`
	MemoryMB     uint64 `json:"memory_mb"`
}

// HealthCheck returns the health status of the application
func HealthCheck(c echo.Context) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			MemoryMB:     bToMb(m.Alloc),
		},
		Checks: map[string]string{
			"database": "not_configured",
			"cache":    "not_configured",
			"storage":  "healthy",
		},
	}

	return c.JSON(http.StatusOK, response)
}

// ReadinessCheck returns whether the application is ready to serve traffic
func ReadinessCheck(c echo.Context) error {
	// Add any readiness checks here (database connections, external services, etc.)
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}

// LivenessCheck returns whether the application is alive
func LivenessCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "alive",
	})
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
