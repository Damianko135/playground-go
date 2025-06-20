package handlers

import (
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

// Metrics holds application metrics
type Metrics struct {
	RequestCount    int64     `json:"request_count"`
	ErrorCount      int64     `json:"error_count"`
	StartTime       time.Time `json:"start_time"`
	LastRequestTime time.Time `json:"last_request_time"`
	Uptime          string    `json:"uptime"`
	GoRoutines      int       `json:"goroutines"`
	MemoryUsage     uint64    `json:"memory_usage_mb"`
	CPUCount        int       `json:"cpu_count"`
}

var (
	requestCount    int64
	errorCount      int64
	lastRequestTime time.Time
	startTime       = time.Now()
)

// IncrementRequestCount increments the request counter
func IncrementRequestCount() {
	atomic.AddInt64(&requestCount, 1)
	lastRequestTime = time.Now()
}

// IncrementErrorCount increments the error counter
func IncrementErrorCount() {
	atomic.AddInt64(&errorCount, 1)
}

// GetMetrics returns application metrics
func GetMetrics(c echo.Context) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := Metrics{
		RequestCount:    atomic.LoadInt64(&requestCount),
		ErrorCount:      atomic.LoadInt64(&errorCount),
		StartTime:       startTime,
		LastRequestTime: lastRequestTime,
		Uptime:          time.Since(startTime).String(),
		GoRoutines:      runtime.NumGoroutine(),
		MemoryUsage:     bToMb(m.Alloc),
		CPUCount:        runtime.NumCPU(),
	}

	return c.JSON(http.StatusOK, metrics)
}

// MetricsMiddleware tracks request metrics
func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			IncrementRequestCount()

			err := next(c)

			if err != nil {
				IncrementErrorCount()
			}

			return err
		}
	}
}
