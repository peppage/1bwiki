package main

import (
	"time"

	"github.com/labstack/echo"
)

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
