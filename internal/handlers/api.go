package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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

// HTMX-specific handlers that return HTML fragments

// GetWeatherHTML returns weather data as HTML fragment for HTMX
func GetWeatherHTML(c echo.Context) error {
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

	html := fmt.Sprintf(`
		<div class="flex justify-between items-center">
			<span class="font-medium">%s</span>
			<span class="text-2xl font-bold text-blue-600">%.1f°C</span>
		</div>
		<p class="text-gray-600">%s</p>
		<div class="flex justify-between text-sm text-gray-500">
			<span>Humidity: %d%%</span>
			<span>Wind: %.1f km/h</span>
		</div>
		<p class="text-xs text-gray-400">Updated: %s</p>
	`, weather.Location, weather.Temperature, weather.Description, weather.Humidity, weather.WindSpeed, weather.Timestamp)

	return c.HTML(http.StatusOK, html)
}

// GetQuoteHTML returns quote data as HTML fragment for HTMX
func GetQuoteHTML(c echo.Context) error {
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
	quote := quotes[index]

	html := fmt.Sprintf(`
		<blockquote class="text-gray-700 italic">"%s"</blockquote>
		<cite class="text-sm text-gray-500 block text-right">— %s</cite>
	`, quote.Text, quote.Author)

	return c.HTML(http.StatusOK, html)
}

// GetColorPaletteHTML returns color palette as HTML fragment for HTMX
func GetColorPaletteHTML(c echo.Context) error {
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
	palette := palettes[index]

	var colorsHTML string
	for _, color := range palette.Colors {
		colorsHTML += fmt.Sprintf(`<div class="w-8 h-8 rounded cursor-pointer hover:scale-110 transition-transform" 
			style="background-color: %s" 
			title="%s" 
			onclick="copyToClipboard('%s')"></div>`, color, color, color)
	}

	html := fmt.Sprintf(`
		<h4 class="font-medium text-gray-900">%s</h4>
		<div class="flex space-x-1">%s</div>
		<p class="text-xs text-gray-500">Theme: %s • Click colors to copy</p>
	`, palette.Name, colorsHTML, palette.Theme)

	return c.HTML(http.StatusOK, html)
}

// GetJokeHTML returns joke as HTML fragment for HTMX
func GetJokeHTML(c echo.Context) error {
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
	joke := jokes[index]

	html := fmt.Sprintf(`
		<p class="text-gray-700 font-medium">%s</p>
		<p class="text-green-600 font-semibold">%s</p>
		<span class="badge">%s</span>
	`, joke.Setup, joke.Punchline, joke.Type)

	return c.HTML(http.StatusOK, html)
}

// GetSystemStatsHTML returns system stats as HTML fragment for HTMX
func GetSystemStatsHTML(c echo.Context) error {
	stats := map[string]interface{}{
		"cpu_usage":    fmt.Sprintf("%.1f%%", 23.5),
		"memory_usage": fmt.Sprintf("%.1f%%", 45.2),
		"disk_usage":   fmt.Sprintf("%.1f%%", 67.8),
		"network_in":   "1.2 MB/s",
		"network_out":  "0.8 MB/s",
		"uptime":       time.Since(startTime).String(),
		"timestamp":    time.Now().Unix(),
	}

	html := fmt.Sprintf(`
		<div class="flex justify-between">
			<span>CPU Usage:</span>
			<span class="font-medium">%s</span>
		</div>
		<div class="flex justify-between">
			<span>Memory:</span>
			<span class="font-medium">%s</span>
		</div>
		<div class="flex justify-between">
			<span>Disk:</span>
			<span class="font-medium">%s</span>
		</div>
		<div class="flex justify-between">
			<span>Uptime:</span>
			<span class="font-medium text-xs">%s</span>
		</div>
	`, stats["cpu_usage"], stats["memory_usage"], stats["disk_usage"], stats["uptime"])

	return c.HTML(http.StatusOK, html)
}

// GetWorldClockHTML returns world clock as HTML fragment for HTMX
func GetWorldClockHTML(c echo.Context) error {
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

	var clocksHTML string
	for city, timeStr := range timeZones {
		clocksHTML += fmt.Sprintf(`
			<div class="text-center p-3 bg-white rounded-lg border border-green-200">
				<div class="text-lg font-bold text-green-600">%s</div>
				<div class="text-sm text-gray-600">%s</div>
			</div>
		`, timeStr, city)
	}

	return c.HTML(http.StatusOK, clocksHTML)
}

// GetRandomNumberHTML returns random number as HTML fragment for HTMX
func GetRandomNumberHTML(c echo.Context) error {
	min := c.QueryParam("min")
	max := c.QueryParam("max")

	minVal := 1
	maxVal := 100

	if min != "" {
		if val, err := strconv.Atoi(min); err == nil {
			minVal = val
		}
	}

	if max != "" {
		if val, err := strconv.Atoi(max); err == nil {
			maxVal = val
		}
	}

	if minVal > maxVal {
		minVal, maxVal = maxVal, minVal
	}

	randomNum := rand.Intn(maxVal-minVal+1) + minVal

	html := fmt.Sprintf(`<span class="text-3xl font-bold text-green-600">%d</span>`, randomNum)

	return c.HTML(http.StatusOK, html)
}
