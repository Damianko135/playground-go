package utils

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Render wraps a templ.Component into an Echo handler
func Render(component templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
}
