package main

import (
	"net/http"
	"strings"
	"time"

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

func Logger() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			if err := h(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := req.Method
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			logger.Debug("", "method", method, "path", path, "code", res.Status(), "time", stop.Sub(start).String())
			return nil
		}
	}
}
