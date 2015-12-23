package main

import (
	"net/http"
	"time"

	mdl "1bwiki/model"
	"1bwiki/tmpl/special"
	"1bwiki/tmpl/special/user"

	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
	"golang.org/x/crypto/bcrypt"
)

func register(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Register(val.(*mdl.User)))
}

func registerHandle(c *echo.Context) error {
	if c.Form("password") == c.Form("passwordConfirm") {
		p, err := bcrypt.GenerateFromPassword([]byte(c.Form("password")), 10)
		if err != nil {
			logger.Error("registering user, encrypting password", "err", err)
		}
		logger.Info("register", "uname", c.Form("username"), "password", p)
		u := mdl.User{
			Name:         c.Form("username"),
			Password:     string(p),
			Registration: time.Now().Unix(),
		}
		err = mdl.CreateUser(&u)
		if err != nil {
			return logger.Error("failed to create user", "err", err)
		}
		session := session.Default(c)
		session.Set("user", u)
		session.Save()
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func login(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Login(val.(*mdl.User)))
}

func loginHandle(c *echo.Context) error {
	u, err := mdl.GetUserByName(c.Form("username"))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized) // The user is invalid!
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Form("password")))
	if err != nil {
		c.Response().Header().Set("Method", "GET")
		return echo.NewHTTPError(http.StatusUnauthorized) // The user is invalid!
	}
	session := session.Default(c)
	session.Set("user", u)
	session.Save()
	return c.Redirect(http.StatusSeeOther, "/")
}

func logout(c *echo.Context) error {
	session := session.Default(c)
	session.Set("user", nil)
	session.Save()
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func prefs(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		return c.HTML(http.StatusOK, user.Prefs(u))
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}
