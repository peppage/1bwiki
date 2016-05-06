package main

import (
	"net"
	"net/http"
	"time"

	mdl "1bwiki/model"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/peppage/echo-middleware/session"
)

func setUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := session.Default(c)
			val := session.Get("user")
			_, ok := val.(*mdl.User)
			if !ok {
				req := c.Request().(*standard.Request).Request
				remoteAddr := req.RemoteAddr
				if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
					remoteAddr = ip
				} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
					remoteAddr = ip
				} else {
					remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
				}
				user := &mdl.User{
					ID:   0,
					Name: remoteAddr,
					Anon: true,
				}
				log.WithFields(log.Fields{
					"user": user.Name,
				}).Warn("Saving anon user")
				session.Set("user", user)
				session.Save()
			}
			return next(c)
		}
	}
}

func serverLogger() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request().(*standard.Request).Request
			res := c.Response().(*standard.Response)

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

			log.WithFields(log.Fields{
				"method": method,
				"path":   path,
				"code":   res.Status(),
				"time":   stop.Sub(start).String(),
			}).Debug()
			return nil
		}
	}
}

func checkLoggedIn() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := session.Default(c)
			val := session.Get("user")
			u, ok := val.(*mdl.User)
			if ok && u.IsLoggedIn() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	}
}

func checkAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := session.Default(c)
			val := session.Get("user")
			u, ok := val.(*mdl.User)
			if ok && u.IsAdmin() {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	}
}
