package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func fixURL() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().Method == "GET" {
				u := strings.Trim(c.Request().URL.String(), "/")
				f := string(u[0])
				if f == strings.ToLower(f) {
					newURL := strings.Replace(strings.ToUpper(f)+string(u[1:]), "%20", "_", -1)
					return c.Redirect(http.StatusMovedPermanently, newURL)
				}
				if strings.Contains(u, "%20") {
					newURL := strings.Replace(u, "%20", "_", -1)
					return c.Redirect(http.StatusMovedPermanently, newURL)
				}
			}
			return next(c)
		}
	}
}
