package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/Damianko135/playground-go/views" // This path must match your `go.mod` module name
)

func main() {
	fmt.Println("🔧 Starting Echo server...")
	e := echo.New()

	// Route to render your templ component
	e.GET("/", func(c echo.Context) error {
		// Render the Init templ component to the Echo response
		return views.Init().Render(c.Request().Context(), c.Response().Writer)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
