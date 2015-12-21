package main

import (
	"encoding/gob"
	"net/http"
	"strings"

	m "1bwiki/model"
	"1bwiki/tmpl/page"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/context"
	"github.com/labstack/echo"
	"github.com/mgutz/logxi/v1"
	"github.com/syntaqx/echo-middleware/session"
)

const secret = "Thisisatemporarysecret"

var logger log.Logger
var store = session.NewCookieStore([]byte(secret))

const noEditArea = "special"

func cleanTitle(t string) string {
	f := string(t[0])
	t = strings.ToUpper(f) + string(t[1:])
	t = strings.Replace(t, "%20", "_", -1)
	t = strings.Replace(t, " ", "_", -1)
	return t
}

func parseTitle(t string) (string, string) {
	URL := strings.Trim(t, "/")
	namespace := ""
	title := ""
	if strings.Contains(URL, ":") {
		split := strings.Split(URL, ":")
		namespace = split[0]
		title = split[1]
	} else {
		title = URL
	}
	return namespace, title
}

func root(c *echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/Main_Page")
}

func wikiPage(c *echo.Context) error {
	n, t := parseTitle(c.Request().URL.String())

	ul := strings.ToLower(c.Request().URL.String())
	if strings.HasPrefix(ul, "/"+noEditArea) {
		return echo.NewHTTPError(http.StatusForbidden, "Editing of special pages disallowed")
	}

	ct := cleanTitle(t)

	if ct != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusMovedPermanently, "/"+n+ct)
	}

	if c.Query("oldid") != "" {
		pv := m.GetOldPageView(c.Query("oldid"))
		session := session.Default(c)
		val := session.Get("user")
		return c.HTML(http.StatusOK, page.Oldversion(val.(*m.User), pv))
	}

	pv := m.GetPageView(n, t)

	if pv.NiceTitle != "" {
		session := session.Default(c)
		val := session.Get("user")
		return c.HTML(http.StatusOK, page.Page(val.(*m.User), pv))
	}
	if n != "" {
		n += ":"
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/special/edit?title="+n+t)
}

func savePage(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*m.User)
	if !ok {
		return logger.Error("User saving page is invalid", "user", u)
	}

	minor := c.Form("minor") == "on"
	p, err := m.CreateOrUpdatePage(u, m.CreatePageOptions{
		Title:     c.Form("title"),
		Namespace: c.Form("namespace"),
		Text:      c.Form("text"),
		Comment:   c.Form("summary"),
		IsMinor:   minor,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Save page  failed")
	}
	return c.Redirect(http.StatusSeeOther, p.Title)
}

func init() {
	gob.Register(&m.User{})
	logger = log.New("1bwiki")
}

func main() {
	m.SetupDb()

	e := echo.New()
	e.Use(session.Sessions("session", store))
	assetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())
	e.Get("/static/*", func(c *echo.Context) error {
		http.StripPrefix("/static/", assetHandler).ServeHTTP(c.Response().Writer(), c.Request())
		return nil
	})
	e.Get("/favicon.ico", func(c *echo.Context) error {
		http.StripPrefix("", assetHandler).ServeHTTP(c.Response().Writer(), c.Request())
		return nil
	})
	e.HTTP2(true)
	e.Use(setUser())
	e.Use(serverLogger())

	e.Get("/", root)
	e.Get("/*", wikiPage)
	e.Post("/save", savePage)

	s := e.Group("/special")
	s.Get("/edit", edit)
	s.Get("/history", history)
	s.Get("/recentchanges", recentChanges)
	s.Get("/register", register)
	s.Post("/register", registerHandle)
	s.Get("/login", login)
	s.Post("/login", loginHandle)
	s.Get("/logout", logout)
	s.Get("/preferences", prefs)
	s.Get("/admin", admin)
	s.Post("/admin", adminHandle)

	http.ListenAndServe(":8000", context.ClearHandler(e))
}
