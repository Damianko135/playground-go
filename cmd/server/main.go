package main

import (
	"fmt"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Damianko135/playground-go/views"
	"github.com/Damianko135/playground-go/internal/utils"
)

func main() {
	fmt.Println("🔧 Starting Echo server...")
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route to render your templ component
	e.GET("/", func(c echo.Context) error {
		// Render the Init templ component to the Echo response
		return views.Home().Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/", utils.Render(views.Home()))
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}



