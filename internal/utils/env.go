package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnvVar returns the value of an environment variable or a fallback value
func GetEnvVar(variable string, fallback string) (string, error) {
	value := os.Getenv(variable)
	if value == "" {
		return fallback, nil
	}
	return value, nil
}

// GetEnvInt returns an environment variable as an integer with a fallback
func GetEnvInt(variable string, fallback int) (int, error) {
	value := os.Getenv(variable)
	if value == "" {
		return fallback, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback, errors.New("environment variable " + variable + " is not a valid integer: " + err.Error())
	}
	return intValue, nil
}

// GetEnvBool returns an environment variable as a boolean with a fallback
// Accepts: true, false, 1, 0, yes, no, on, off (case insensitive)
func GetEnvBool(variable string, fallback bool) (bool, error) {
	value := os.Getenv(variable)
	if value == "" {
		return fallback, nil
	}

	boolValue, err := parseBool(value)
	if err != nil {
		return fallback, errors.New("environment variable " + variable + " is not a valid boolean: " + err.Error())
	}
	return boolValue, nil
}

// GetEnvDuration returns an environment variable as a time.Duration with a fallback
// Accepts values like "10s", "5m", "1h", "300ms", etc.
func GetEnvDuration(variable string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(variable)
	if value == "" {
		return fallback, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback, errors.New("environment variable " + variable + " is not a valid duration: " + err.Error())
	}
	return duration, nil
}

// Helper function to parse boolean values more flexibly
func parseBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, errors.New("invalid boolean value: " + value)
	}
}
