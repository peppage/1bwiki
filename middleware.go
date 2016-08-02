package main

import (
	"net/http"

	mdl "1bwiki/model"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
)

type sessionMiddleware struct{}

func (s *sessionMiddleware) Serve(c *iris.Context) {
	val := c.Session().Get("user")
	_, ok := val.(*mdl.User)
	if !ok {
		remoteAddr := c.RemoteAddr()
		user := &mdl.User{
			ID:       0,
			Name:     remoteAddr,
			Anon:     true,
			TimeZone: "UTC",
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
