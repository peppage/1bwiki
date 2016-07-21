package main

import (
	"net/http"
	"time"

	mdl "1bwiki/model"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

type sessionMiddleware struct{}

func (s *sessionMiddleware) Serve(c *iris.Context) {
	val := c.Session().Get("user")
	_, ok := val.(*mdl.User)
	if !ok {
		remoteAddr := c.RemoteAddr()
		user := &mdl.User{
			ID:   0,
			Name: remoteAddr,
			Anon: true,
		}
		log.WithFields(log.Fields{
			"user": user.Name,
		}).Warn("saving anon user")
		c.Session().Set("user", user)
	}
	c.Next()
}

type loggedInMiddleware struct{}

func (s *loggedInMiddleware) Serve(c *iris.Context) {
	val := c.Session().Get("user")
	u, ok := val.(*mdl.User)
	if ok && u.IsLoggedIn() {
		c.Next()
	}
	c.Error("Must be logged in ", http.StatusUnauthorized)
	return
}

type adminMiddleware struct{}

func (s *adminMiddleware) Serve(c *iris.Context) {
	val := c.Session().Get("user")
	u, ok := val.(*mdl.User)
	if ok && u.IsAdmin() {
		c.Next()
	}
	c.Error("Must be admin ", http.StatusUnauthorized)
	return
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
