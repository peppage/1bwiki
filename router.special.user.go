package main

import (
	"net/http"
	"time"

	m "1bwiki/model"
	"1bwiki/tmpl/special"

	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
	"golang.org/x/crypto/bcrypt"
)

func register(c *echo.Context) error {
	return c.HTML(http.StatusOK, special.Register())
}

func registerHandle(c *echo.Context) error {
	if c.Form("password") == c.Form("passwordConfirm") {
		p, err := bcrypt.GenerateFromPassword([]byte(c.Form("password")), 10)
		if err != nil {
			logger.Error("registering user, encrypting password", "err", err)
		}
		logger.Info("register", "uname", c.Form("username"), "password", p)
		u := m.User{
			Name:         c.Form("username"),
			Password:     string(p),
			Registration: time.Now().Unix(),
		}
		u.Create()
	}

	return c.String(http.StatusOK, "ok")
}

func login(c *echo.Context) error {
	return c.HTML(http.StatusOK, special.Login())
}

func loginHandle(c *echo.Context) error {
	u := m.User{
		Name:     c.Form("username"),
		Password: c.Form("password"),
	}
	err := u.Verify()
	if err != nil {
		c.Response().Header().Set("Method", "GET")
		return echo.NewHTTPError(http.StatusUnauthorized) // The user is invalid!
	}
	session := session.Default(c)
	session.Set("user", u)
	session.Save()
	return c.Redirect(http.StatusSeeOther, "/admin")
}

func logout(c *echo.Context) error {
	session := session.Default(c)
	session.Set("user", nil)
	session.Save()
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
