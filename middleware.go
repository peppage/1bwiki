package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func fixUrl() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().Method == "GET" {
				u := strings.Trim(c.Request().URL.String(), "/")
				f := string(u[0])
				if f == strings.ToLower(f) {
					newUrl := strings.Replace(strings.ToUpper(f)+string(u[1:]), "%20", "_", -1)
					return c.Redirect(http.StatusMovedPermanently, newUrl)
				}
				if strings.Contains(u, "%20") {
					newUrl := strings.Replace(u, "%20", "_", -1)
					return c.Redirect(http.StatusMovedPermanently, newUrl)
				}
			}
			return next(c)
		}
	}
}
