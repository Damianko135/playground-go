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

// GetRequiredVar returns the value of an environment variable or an error if not set
func GetRequiredVar(variable string) (string, error) {
	value := os.Getenv(variable)
	if value == "" {
		return "", errors.New("required environment variable " + variable + " is not set")
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

// GetRequiredInt returns an environment variable as an integer or an error if not set or invalid
func GetRequiredInt(variable string) (int, error) {
	value := os.Getenv(variable)
	if value == "" {
		return 0, errors.New("required environment variable " + variable + " is not set")
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("environment variable " + variable + " is not a valid integer: " + err.Error())
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

// GetRequiredBool returns an environment variable as a boolean or an error if not set or invalid
func GetRequiredBool(variable string) (bool, error) {
	value := os.Getenv(variable)
	if value == "" {
		return false, errors.New("required environment variable " + variable + " is not set")
	}

	boolValue, err := parseBool(value)
	if err != nil {
		return false, errors.New("environment variable " + variable + " is not a valid boolean: " + err.Error())
	}
	return boolValue, nil
}

// GetEnvFloat returns an environment variable as a float64 with a fallback
func GetEnvFloat(variable string, fallback float64) (float64, error) {
	value := os.Getenv(variable)
	if value == "" {
		return fallback, nil
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback, errors.New("environment variable " + variable + " is not a valid float: " + err.Error())
	}
	return floatValue, nil
}

// GetRequiredFloat returns an environment variable as a float64 or an error if not set or invalid
func GetRequiredFloat(variable string) (float64, error) {
	value := os.Getenv(variable)
	if value == "" {
		return 0, errors.New("required environment variable " + variable + " is not set")
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.New("environment variable " + variable + " is not a valid float: " + err.Error())
	}
	return floatValue, nil
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

// GetRequiredDuration returns an environment variable as a time.Duration or an error if not set or invalid
func GetRequiredDuration(variable string) (time.Duration, error) {
	value := os.Getenv(variable)
	if value == "" {
		return 0, errors.New("required environment variable " + variable + " is not set")
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, errors.New("environment variable " + variable + " is not a valid duration: " + err.Error())
	}
	return duration, nil
}

// GetEnvSlice returns an environment variable as a slice of strings, split by separator
// Default separator is comma ","
func GetEnvSlice(variable string, separator string, fallback []string) ([]string, error) {
	value := os.Getenv(variable)
	if value == "" {
		return fallback, nil
	}

	if separator == "" {
		separator = ","
	}

	parts := strings.Split(value, separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result, nil
}

// GetRequiredSlice returns an environment variable as a slice of strings or an error if not set
func GetRequiredSlice(variable string, separator string) ([]string, error) {
	value := os.Getenv(variable)
	if value == "" {
		return nil, errors.New("required environment variable " + variable + " is not set")
	}

	if separator == "" {
		separator = ","
	}

	parts := strings.Split(value, separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("environment variable " + variable + " is empty or contains only separators")
	}

	return result, nil
}

// IsEnvSet checks if an environment variable is set (even if empty)
func IsEnvSet(variable string) bool {
	_, exists := os.LookupEnv(variable)
	return exists
}

// GetEnvWithValidation returns an environment variable with custom validation
func GetEnvWithValidation(variable string, fallback string, validator func(string) error) (string, error) {
	value := os.Getenv(variable)
	if value == "" {
		if validator != nil {
			if err := validator(fallback); err != nil {
				return "", errors.New("fallback value for " + variable + " failed validation: " + err.Error())
			}
		}
		return fallback, nil
	}

	if validator != nil {
		if err := validator(value); err != nil {
			return "", errors.New("environment variable " + variable + " failed validation: " + err.Error())
		}
	}

	return value, nil
}

// MustGetEnv returns an environment variable or panics if not set
func MustGetEnv(variable string) string {
	value, err := GetRequiredVar(variable)
	if err != nil {
		panic(err)
	}
	return value
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
