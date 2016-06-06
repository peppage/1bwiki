package main

import (
	"net/http"

	mdl "1bwiki/model"
	"1bwiki/view"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/peppage/echo-middleware/session"
)

func admin(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u := val.(*mdl.User)
	p := &view.AdminPage{
		User: u,
		URL:  "/special/admin",
	}
	return c.HTML(http.StatusOK, view.PageTemplate(p))
}

func adminHandle(c echo.Context) error {
	err := mdl.SetAnonEditing(c.FormValue("anon") == "on")
	if err != nil {
		log.WithError(err).Error("Failed to set anon editing")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	err = mdl.SetSignups(c.FormValue("signup") == "on")
	if err != nil {
		log.WithError(err).Error("Failed to set signups")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Redirect(http.StatusSeeOther, "/special/admin")
}
