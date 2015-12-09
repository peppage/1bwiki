package main

import (
	"net"
	"time"

	m "1bwiki/model"

	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
)

func setUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			session := session.Default(c)
			val := session.Get("user")
			_, ok := val.(*m.User)
			if !ok {
				req := c.Request()
				remoteAddr := req.RemoteAddr
				if ip := req.Header.Get(echo.XRealIP); ip != "" {
					remoteAddr = ip
				} else if ip = req.Header.Get(echo.XForwardedFor); ip != "" {
					remoteAddr = ip
				} else {
					remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
				}
				user := &m.User{
					ID:   0,
					Name: remoteAddr,
					Anon: true,
				}
				logger.Warn("Saving anon user", "user", user)
				session.Set("user", user)
				session.Save()
			}
			return next(c)
		}
	}
}

func serverLogger() echo.MiddlewareFunc {
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
