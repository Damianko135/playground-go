package main

import (
	"fmt"
	"os"

	"github.com/Damianko135/playground-go/internal/utils"
	"github.com/Damianko135/playground-go/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("🔧 Starting Echo server...")
	e := echo.New()

	// Enable debug mode and logging in development
	if os.Getenv("GO_ENV") == "development" {
		e.Debug = true
		// Uncomment below for request logging during development
		e.Use(middleware.Logger())
		fmt.Println("🐛 Debug mode enabled")
	}
	e.Use(middleware.Recover())

	e.Static("/static", "./static")
	e.Static("/", "./public")

	// Serve homepage using your render helper
	e.GET("/", utils.Temple(views.Home()))
	e.GET("/about", utils.Temple(views.About()))

	// Use PORT environment variable or default to 8080
	port, err := utils.GetEnvVar("PORT", "8080")
	if err != nil {
		e.Logger.Fatal(fmt.Sprintf("Failed to get PORT: %v", err))
		return
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
