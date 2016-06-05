package main

import (
	"net/http"
	"strings"
	"time"

	mdl "1bwiki/model"
	"1bwiki/tmpl/special"
	"1bwiki/tmpl/special/user"
	"1bwiki/view"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/peppage/echo-middleware/session"
)

func register(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	flashes := session.Flashes()
	session.Save()
	p := &view.RegisterPage{
		User:     val.(*mdl.User),
		URL:      "/special/register",
		Messages: flashes,
	}
	return c.HTML(http.StatusOK, view.PageTemplate(p))
}

func registerHandle(c echo.Context) error {
	session := session.Default(c)
	if c.FormValue("password") == c.FormValue("passwordConfirm") {
		u := mdl.User{
			Name:         c.FormValue("username"),
			Password:     c.FormValue("password"),
			Registration: time.Now().Unix(),
		}
		err := mdl.CreateUser(&u)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				session.AddFlash("Username already exists")
				session.Save()
				return c.Redirect(http.StatusSeeOther, "/special/register")
			}
			log.WithFields(log.Fields{
				"err":  err,
				"user": c.FormValue("username"),
			})
			return echo.NewHTTPError(http.StatusBadRequest, "failed to create user")
		}
		session.Set("user", u)
		session.Save()
	} else {
		session.AddFlash("Passwords don't match")
		session.Save()
		return c.Redirect(http.StatusSeeOther, "/special/register")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func login(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Login(val.(*mdl.User)))
}

func loginHandle(c echo.Context) error {
	u, err := mdl.GetUserByName(c.FormValue("username"))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized) // The user is doesn't exist
	}

	if u.ValidatePassword(c.FormValue("password")) {
		session := session.Default(c)
		session.Set("user", u)
		session.Save()
		return c.Redirect(http.StatusSeeOther, "/")
	}

	c.Response().Header().Set("Method", "GET")
	return echo.NewHTTPError(http.StatusUnauthorized)
}

func logout(c echo.Context) error {
	session := session.Default(c)
	session.Set("user", nil)
	session.Save()
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func prefs(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		return c.HTML(http.StatusOK, user.Prefs(u))
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}

func prefsPasword(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		return c.HTML(http.StatusOK, user.Password(u))
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}

func handlePrefsPassword(c echo.Context) error {
	if c.FormValue("newpassword1") != c.FormValue("newpassword2") {
		// need to implement better
		return echo.NewHTTPError(http.StatusBadRequest, "password do not match")
	}
	session := session.Default(c)
	val := session.Get("user")
	u, _ := val.(*mdl.User)

	if u.ValidatePassword(c.FormValue("oldpassword")) {
		u.Password = c.FormValue("newpassword1")
		err := u.EncodePassword()
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"user": u.Name,
			})
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed.")
		}
		err = mdl.UpdateUser(u)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.Redirect(http.StatusSeeOther, "/special/preferences")
	}

	c.Response().Header().Set("Method", "GET")
	return echo.NewHTTPError(http.StatusUnauthorized) // The user is invalid!
}
