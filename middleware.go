package main

import (
	"net/http"

	mdl "1bwiki/model"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

type sessionMiddleware struct{}

func (s *sessionMiddleware) Serve(c *iris.Context) {
	val := c.Session().Get("user")
	_, ok := val.(*mdl.User)
	if !ok {
		remoteAddr := c.RemoteAddr()
		user := &mdl.User{
			ID:         0,
			Name:       remoteAddr,
			Anon:       true,
			TimeZone:   "UTC",
			DateFormat: "15:04, 2 January 2006",
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
		return
	}
	c.Redirect("/special/login", http.StatusUnauthorized)
}

type adminMiddleware struct{}

func (s *adminMiddleware) Serve(c *iris.Context) {
	val := c.Session().Get("user")
	u, ok := val.(*mdl.User)
	if ok && u.IsAdmin() {
		c.Next()
		return
	}
	c.Redirect("/special/login", http.StatusUnauthorized)
}
