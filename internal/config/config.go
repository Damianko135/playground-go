package config

import (
	"errors"
	"strconv"
	"time"

	"github.com/Damianko135/playground-go/internal/utils"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	API      APIConfig
	Features FeatureConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	Host         string
	Environment  string
	Debug        bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// APIConfig holds API-related configuration
type APIConfig struct {
	Key        string
	RateLimit  int
	EnableCORS bool
	EnableGzip bool
}

// FeatureConfig holds feature flags
type FeatureConfig struct {
	EnableHealthCheck bool
	EnableMetrics     bool
	EnableProfiling   bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port, err := utils.GetEnvVar("PORT", "8080")
	if err != nil {
		return nil, err
	}

	host, err := utils.GetEnvVar("HOST", "localhost")
	if err != nil {
		return nil, err
	}

	environment, err := utils.GetEnvVar("GO_ENV", "development")
	if err != nil {
		return nil, err
	}

	debug, err := utils.GetEnvBool("DEBUG", environment == "development")
	if err != nil {
		return nil, err
	}

	readTimeout, err := utils.GetEnvDuration("READ_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, err
	}

	writeTimeout, err := utils.GetEnvDuration("WRITE_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, err
	}

	apiKey, err := utils.GetEnvVar("API_KEY", "")
	if err != nil {
		return nil, err
	}

	rateLimit, err := utils.GetEnvInt("RATE_LIMIT", 20)
	if err != nil {
		return nil, err
	}

	enableCORS, err := utils.GetEnvBool("ENABLE_CORS", true)
	if err != nil {
		return nil, err
	}

	enableGzip, err := utils.GetEnvBool("ENABLE_GZIP", true)
	if err != nil {
		return nil, err
	}

	enableHealthCheck, err := utils.GetEnvBool("ENABLE_HEALTH_CHECK", true)
	if err != nil {
		return nil, err
	}

	enableMetrics, err := utils.GetEnvBool("ENABLE_METRICS", false)
	if err != nil {
		return nil, err
	}

	enableProfiling, err := utils.GetEnvBool("ENABLE_PROFILING", debug)
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port:         port,
			Host:         host,
			Environment:  environment,
			Debug:        debug,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		API: APIConfig{
			Key:        apiKey,
			RateLimit:  rateLimit,
			EnableCORS: enableCORS,
			EnableGzip: enableGzip,
		},
		Features: FeatureConfig{
			EnableHealthCheck: enableHealthCheck,
			EnableMetrics:     enableMetrics,
			EnableProfiling:   enableProfiling,
		},
	}, nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetAddress returns the full server address
func (c *Config) GetAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// GetPort returns just the port number
func (c *Config) GetPort() string {
	return c.Server.Port
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate port
	if port, err := strconv.Atoi(c.Server.Port); err != nil || port < 1 || port > 65535 {
		return err
	}

	// Validate timeouts
	if c.Server.ReadTimeout < 0 || c.Server.WriteTimeout < 0 {
		return errors.New("server timeouts must be positive")
	}

	// Validate rate limit
	if c.API.RateLimit < 1 {
		return errors.New("rate limit must be at least 1")
	}

	return nil
}

// Print prints the configuration (without sensitive data)
func (c *Config) Print() {
	println("🔧 Configuration:")
	println("  Server:")
	println("    Port:", c.Server.Port)
	println("    Host:", c.Server.Host)
	println("    Environment:", c.Server.Environment)
	println("    Debug:", c.Server.Debug)
	println("    Read Timeout:", c.Server.ReadTimeout.String())
	println("    Write Timeout:", c.Server.WriteTimeout.String())
	println("  API:")
	println("    Rate Limit:", c.API.RateLimit)
	println("    Enable CORS:", c.API.EnableCORS)
	println("    Enable Gzip:", c.API.EnableGzip)
	if c.API.Key != "" {
		println("    API Key: [CONFIGURED]")
	} else {
		println("    API Key: [NOT SET]")
	}
	println("  Features:")
	println("    Health Check:", c.Features.EnableHealthCheck)
	println("    Metrics:", c.Features.EnableMetrics)
	println("    Profiling:", c.Features.EnableProfiling)
}
