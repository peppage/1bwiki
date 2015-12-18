package main

import (
	"net/http"

	m "1bwiki/model"
	"1bwiki/tmpl/special"

	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
)

func admin(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u := val.(*m.User)
	if u.IsAdmin() {
		return c.HTML(http.StatusOK, special.Admin(u, m.AnonEditing(), m.Signups()))
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "Admin only")
}

func adminHandle(c *echo.Context) error {
	err := m.SetAnonEditing(c.Form("anon") == "on")
	if err != nil {
		logger.Error("admin handler", "set anon db", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	err = m.SetSignups(c.Form("signup") == "on")
	if err != nil {
		logger.Error("admin handler", "set signup db", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Redirect(http.StatusSeeOther, "/special/admin")
}
