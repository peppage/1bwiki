package main

import (
	"net/http"
	"strings"
	"time"

	m "1bwiki/model"
	"1bwiki/tmpl"
	"1bwiki/tmpl/special"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func action(c *echo.Context) error {
	n, t := parseTitle(c.Query("title"))
	if c.Query("action") == "edit" {
		pv := m.GetPageView(n, t)
		if pv.NiceTitle == "" {
			if n != "" {
				n += ":"
			}
			pv.NameSpace = n
			pv.Title = t
			pv.NiceTitle = strings.Replace(n+t, "_", " ", -1)
		}
		return c.HTML(http.StatusOK, tmpl.Editpage(pv))
	} else if c.Query("action") == "history" {
		revs, err := m.GetPageRevisions(c.Query("title"))
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, "")
		}
		return c.HTML(http.StatusOK, tmpl.Pagerevs(c.Query("title"), revs))
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Not an acceptable action")
}

func recentChanges(c *echo.Context) error {
	revs, err := m.GetRevisions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	return c.HTML(http.StatusOK, tmpl.Recentchanges(revs))
}

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
