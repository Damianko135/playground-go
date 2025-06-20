package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// WeatherData represents weather information
type WeatherData struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Icon        string  `json:"icon"`
	Timestamp   string  `json:"timestamp"`
}

// QuoteData represents an inspirational quote
type QuoteData struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

// GetWeather returns mock weather data (in a real app, you'd call a weather API)
func GetWeather(c echo.Context) error {
	// Mock weather data - in production, you'd call a real weather API
	weather := WeatherData{
		Location:    "Amsterdam, NL",
		Temperature: 18.5,
		Description: "Partly cloudy",
		Humidity:    65,
		WindSpeed:   12.3,
		Icon:        "partly-cloudy",
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return c.JSON(http.StatusOK, weather)
}

// GetQuote returns a random inspirational quote
func GetQuote(c echo.Context) error {
	quotes := []QuoteData{
		{Text: "The only way to do great work is to love what you do.", Author: "Steve Jobs"},
		{Text: "Innovation distinguishes between a leader and a follower.", Author: "Steve Jobs"},
		{Text: "Code is like humor. When you have to explain it, it's bad.", Author: "Cory House"},
		{Text: "First, solve the problem. Then, write the code.", Author: "John Johnson"},
		{Text: "Experience is the name everyone gives to their mistakes.", Author: "Oscar Wilde"},
		{Text: "In order to be irreplaceable, one must always be different.", Author: "Coco Chanel"},
		{Text: "Java is to JavaScript what car is to Carpet.", Author: "Chris Heilmann"},
		{Text: "Knowledge is power.", Author: "Francis Bacon"},
		{Text: "Sometimes it pays to stay in bed on Monday, rather than spending the rest of the week debugging Monday's code.", Author: "Dan Salomon"},
		{Text: "Perfection is achieved not when there is nothing more to add, but rather when there is nothing more to take away.", Author: "Antoine de Saint-Exupery"},
	}

	// Simple random selection based on current time
	index := int(time.Now().Unix()) % len(quotes)
	return c.JSON(http.StatusOK, quotes[index])
}

// GetSystemStats returns system performance statistics
func GetSystemStats(c echo.Context) error {
	stats := map[string]interface{}{
		"cpu_usage":    fmt.Sprintf("%.1f%%", 23.5),
		"memory_usage": fmt.Sprintf("%.1f%%", 45.2),
		"disk_usage":   fmt.Sprintf("%.1f%%", 67.8),
		"network_in":   "1.2 MB/s",
		"network_out":  "0.8 MB/s",
		"uptime":       time.Since(startTime).String(),
		"timestamp":    time.Now().Unix(),
	}

	return c.JSON(http.StatusOK, stats)
}

// ColorPalette represents a color palette
type ColorPalette struct {
	Name   string   `json:"name"`
	Colors []string `json:"colors"`
	Theme  string   `json:"theme"`
}

// GetColorPalette returns a random color palette
func GetColorPalette(c echo.Context) error {
	palettes := []ColorPalette{
		{
			Name:   "Ocean Breeze",
			Colors: []string{"#0077be", "#00a8cc", "#40e0d0", "#87ceeb", "#b0e0e6"},
			Theme:  "cool",
		},
		{
			Name:   "Sunset Glow",
			Colors: []string{"#ff6b35", "#f7931e", "#ffd700", "#ff69b4", "#ff1493"},
			Theme:  "warm",
		},
		{
			Name:   "Forest Calm",
			Colors: []string{"#228b22", "#32cd32", "#90ee90", "#98fb98", "#f0fff0"},
			Theme:  "nature",
		},
		{
			Name:   "Purple Dreams",
			Colors: []string{"#4b0082", "#8a2be2", "#9370db", "#ba55d3", "#dda0dd"},
			Theme:  "mystical",
		},
		{
			Name:   "Monochrome",
			Colors: []string{"#000000", "#404040", "#808080", "#c0c0c0", "#ffffff"},
			Theme:  "neutral",
		},
	}

	index := int(time.Now().Unix()) % len(palettes)
	return c.JSON(http.StatusOK, palettes[index])
}

// JokeData represents a programming joke
type JokeData struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
	Type      string `json:"type"`
}

// GetJoke returns a random programming joke
func GetJoke(c echo.Context) error {
	jokes := []JokeData{
		{
			Setup:     "Why do programmers prefer dark mode?",
			Punchline: "Because light attracts bugs!",
			Type:      "programming",
		},
		{
			Setup:     "How many programmers does it take to change a light bulb?",
			Punchline: "None. That's a hardware problem.",
			Type:      "programming",
		},
		{
			Setup:     "Why do Java developers wear glasses?",
			Punchline: "Because they can't C#!",
			Type:      "programming",
		},
		{
			Setup:     "What's the object-oriented way to become wealthy?",
			Punchline: "Inheritance.",
			Type:      "programming",
		},
		{
			Setup:     "Why did the programmer quit his job?",
			Punchline: "He didn't get arrays.",
			Type:      "programming",
		},
	}

	index := int(time.Now().Unix()) % len(jokes)
	return c.JSON(http.StatusOK, jokes[index])
}

// GetRandomNumber generates a random number within specified range
func GetRandomNumber(c echo.Context) error {
	min := 1
	max := 100

	// Get query parameters if provided
	if minParam := c.QueryParam("min"); minParam != "" {
		if parsedMin, err := parseIntParam(minParam); err == nil {
			min = parsedMin
		}
	}
	if maxParam := c.QueryParam("max"); maxParam != "" {
		if parsedMax, err := parseIntParam(maxParam); err == nil {
			max = parsedMax
		}
	}

	// Ensure min <= max
	if min > max {
		min, max = max, min
	}

	// Generate random number using current time as seed
	randomNum := min + int(time.Now().UnixNano())%(max-min+1)

	result := map[string]interface{}{
		"number":    randomNum,
		"min":       min,
		"max":       max,
		"timestamp": time.Now().Unix(),
	}

	return c.JSON(http.StatusOK, result)
}

// Helper function to parse integer parameters
func parseIntParam(param string) (int, error) {
	var result int
	_, err := fmt.Sscanf(param, "%d", &result)
	return result, err
}

// GetTimeZones returns current time in different time zones
func GetTimeZones(c echo.Context) error {
	now := time.Now()

	timeZones := map[string]string{
		"UTC":         now.UTC().Format("15:04:05"),
		"New York":    now.In(mustLoadLocation("America/New_York")).Format("15:04:05"),
		"London":      now.In(mustLoadLocation("Europe/London")).Format("15:04:05"),
		"Tokyo":       now.In(mustLoadLocation("Asia/Tokyo")).Format("15:04:05"),
		"Sydney":      now.In(mustLoadLocation("Australia/Sydney")).Format("15:04:05"),
		"Los Angeles": now.In(mustLoadLocation("America/Los_Angeles")).Format("15:04:05"),
		"Amsterdam":   now.In(mustLoadLocation("Europe/Amsterdam")).Format("15:04:05"),
		"Singapore":   now.In(mustLoadLocation("Asia/Singapore")).Format("15:04:05"),
	}

	result := map[string]interface{}{
		"timezones": timeZones,
		"timestamp": now.Unix(),
		"iso_time":  now.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, result)
}

// Helper function to load time zone location
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}
