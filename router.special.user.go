package main

import (
	"net/http"
	"strings"
	"time"

	mdl "1bwiki/model"
	"1bwiki/view"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
)

func register(c *iris.Context) {
	val := c.Session().Get("user")
	flash, _ := c.GetFlash("error")
	p := &view.RegisterPage{
		User:    val.(*mdl.User),
		URL:     "/special/register",
		Message: flash,
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func registerHandle(c *iris.Context) {
	if c.FormValueString("password") == c.FormValueString("passwordConfirm") {
		u := mdl.User{
			Name:         c.FormValueString("username"),
			Password:     c.FormValueString("password"),
			Registration: time.Now().Unix(),
		}
		err := mdl.CreateUser(&u)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				c.SetFlash("error", "Username already exists")
				c.Redirect("/special/register", http.StatusSeeOther)
				return
			}
			log.WithFields(log.Fields{
				"err":  err,
				"user": c.FormValue("username"),
			})
			c.Error("Failed to Register", http.StatusBadRequest)
			return
		}
		c.Session().Set("user", u)
	} else {
		c.SetFlash("error", "Passwords don't match")
		c.Redirect("/special/register", http.StatusSeeOther)
		return
	}

	c.Redirect("/", http.StatusSeeOther)
}

func login(c *iris.Context) {
	val := c.Session().Get("user")
	p := &view.LoginPage{
		User: val.(*mdl.User),
		URL:  "/special/login",
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func loginHandle(c *iris.Context) {
	u, err := mdl.GetUserByName(c.FormValueString("username"))
	if err != nil {
		c.Error("User doesn't exist? Maybe.", http.StatusUnauthorized)
		return
	}

	if u.ValidatePassword(c.FormValueString("password")) {
		c.Session().Set("user", u)
		c.Redirect("/", http.StatusSeeOther)
		return
	}

	c.Error("Failed login", http.StatusUnauthorized)
}

func logout(c *iris.Context) {
	c.Session().Set("user", nil)
	c.Redirect("/", http.StatusTemporaryRedirect)
}

func prefs(c *iris.Context) {
	val := c.Session().Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		p := &view.PrefsPage{
			User: u,
			URL:  "/special/preferences",
		}
		view.WritePageTemplate(c.GetRequestCtx(), p)
		c.HTML(http.StatusOK, "")
	}
}

func handlePrefs(c *iris.Context) {
	val := c.Session().Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		u.TimeZone = c.FormValueString("timezone")
		err := mdl.UpdateUserSettings(u)
		if err != nil {
			c.Error("Failed saving user", http.StatusInternalServerError)
			return
		}
		c.Redirect("/special/preferences", http.StatusSeeOther)
	}
}

func prefsPasword(c *iris.Context) {
	val := c.Session().Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		p := &view.PasswordPage{
			User: u,
			URL:  "/special/preferences/password",
		}
		view.WritePageTemplate(c.GetRequestCtx(), p)
		c.HTML(http.StatusOK, "")
	}
}

func handlePrefsPassword(c *iris.Context) {
	if c.FormValueString("newpassword1") != c.FormValueString("newpassword2") {
		// need to implement better
		c.Error("Passwords do not match", http.StatusBadRequest)
		return
	}
	val := c.Session().Get("user")
	u, _ := val.(*mdl.User)

	if u.ValidatePassword(c.FormValueString("oldpassword")) {
		u.Password = c.FormValueString("newpassword1")
		err := u.EncodePassword()
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"user": u.Name,
			})
			c.Error("Failed", http.StatusInternalServerError)
			return
		}
		err = mdl.UpdateUserPassword(u)
		if err != nil {
			c.EmitError(http.StatusInternalServerError)
			return
		}
		c.Redirect("/special/preferences", http.StatusSeeOther)
	}
}
