package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomLogger returns a custom logger middleware
func CustomLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} | ${status} | ${latency_human} | ${method} ${uri}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
		Output:           nil, // Use default output
	})
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; script-src 'self' 'unsafe-inline'; img-src 'self' data: https:;")
			return next(c)
		}
	}
}

// RateLimiter implements a simple rate limiting middleware
func RateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)) // 20 requests per second
}

// RequestID adds a unique request ID to each request
func RequestID() echo.MiddlewareFunc {
	return middleware.RequestID()
}

// CORS configures Cross-Origin Resource Sharing
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})
}

// Gzip enables gzip compression
func Gzip() echo.MiddlewareFunc {
	return middleware.Gzip()
}

// Cache adds cache headers for static content
func Cache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Add cache headers for static assets
			if isStaticAsset(c.Request().URL.Path) {
				c.Response().Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
				c.Response().Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(time.RFC1123))
			} else {
				// For dynamic content, prevent caching
				c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				c.Response().Header().Set("Pragma", "no-cache")
				c.Response().Header().Set("Expires", "0")
			}
			return next(c)
		}
	}
}

// isStaticAsset checks if the path is for a static asset
func isStaticAsset(path string) bool {
	staticPaths := []string{"/static/", "/css/", "/js/", "/images/", "/fonts/"}
	for _, staticPath := range staticPaths {
		if len(path) >= len(staticPath) && path[:len(staticPath)] == staticPath {
			return true
		}
	}
	return false
}

// APIKeyAuth provides simple API key authentication for API endpoints
func APIKeyAuth(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip auth for non-API endpoints
			if !isAPIEndpoint(c.Request().URL.Path) {
				return next(c)
			}

			// Check for API key in header or query parameter
			key := c.Request().Header.Get("X-API-Key")
			if key == "" {
				key = c.QueryParam("api_key")
			}

			// For demo purposes, we'll allow requests without API key
			// In production, you'd want to enforce this
			if apiKey != "" && key != apiKey {
				return echo.NewHTTPError(401, "Invalid API key")
			}

			return next(c)
		}
	}
}

// isAPIEndpoint checks if the path is an API endpoint
func isAPIEndpoint(path string) bool {
	return len(path) >= 4 && path[:4] == "/api"
}

// ResponseTime adds response time header
func ResponseTime() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)
			c.Response().Header().Set("X-Response-Time", fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1e6))
			return err
		}
	}
}
