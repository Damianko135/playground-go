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
		// e.Use(middleware.Logger())
		fmt.Println("🐛 Debug mode enabled")
	}

	e.Use(middleware.Recover())

	// Serve homepage using your render helper
	e.GET("/", utils.Render(views.Home()))

	// Use PORT environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
