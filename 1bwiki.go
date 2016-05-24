package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"

	mdl "1bwiki/model"
	"1bwiki/setting"
	"1bwiki/tmpl/page"
	"1bwiki/view"

	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/peppage/echo-middleware/session"
)

var store session.CookieStore

const noEditArea = "special"

func convertTitleToUrl(t string) string {
	firstChar := string(t[0])
	t = strings.ToUpper(firstChar) + string(t[1:])
	t = strings.Replace(t, "%20", "_", -1)
	t = strings.Replace(t, " ", "_", -1)
	return t
}

func seperateNamespaceAndTitle(t string) (namespace string, title string) {
	URL := strings.Trim(t, "/")
	if strings.Contains(URL, ":") {
		split := strings.Split(URL, ":")
		namespace = split[0]
		title = split[1]
	} else {
		title = URL
	}
	return namespace, title
}

func root(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/Main_Page")
}

func wikiPage(c echo.Context) error {
	n, t := seperateNamespaceAndTitle(c.Request().(*standard.Request).Request.URL.String())

	ul := strings.ToLower(c.Request().(*standard.Request).Request.URL.String())
	if strings.HasPrefix(ul, "/"+noEditArea) {
		return echo.NewHTTPError(http.StatusForbidden, "Editing of special pages disallowed")
	}

	urlTitle := convertTitleToUrl(t)

	if urlTitle != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusMovedPermanently, "/"+n+urlTitle)
	}

	if c.QueryParam("oldid") != "" {
		pv, err := mdl.GetPageVeiwByID(c.QueryParam("oldid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		session := session.Default(c)
		val := session.Get("user")
		if c.QueryParam("diff") != "" {
			pv2, err := mdl.GetPageVeiwByID(c.QueryParam("diff"))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			p := &view.ArticleDiff{
				User:  val.(*mdl.User),
				Page:  pv,
				Page2: pv2,
			}
			return c.HTML(http.StatusOK, view.PageTemplate(p))
		}
		return c.HTML(http.StatusOK, page.Oldversion(val.(*mdl.User), pv))
	}

	pv := mdl.GetPageView(n, t)

	if pv.NiceTitle != "" && !pv.Deleted {
		session := session.Default(c)
		val := session.Get("user")
		p := &view.Article{
			User: val.(*mdl.User),
			Page: pv,
		}
		return c.HTML(http.StatusOK, view.PageTemplate(p))
	}
	if n != "" {
		n += ":"
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/special/edit?title="+n+t)
}

func savePage(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*mdl.User)
	if !ok {
		log.WithFields(log.Fields{
			"user": u,
		}).Error("User saving page is invalid")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid User")
	}

	minor := c.FormValue("minor") == "on"
	p, err := mdl.CreateOrUpdatePage(u, mdl.CreatePageOptions{
		Title:     c.FormValue("title"),
		Namespace: c.FormValue("namespace"),
		Text:      c.FormValue("text"),
		Comment:   c.FormValue("summary"),
		IsMinor:   minor,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Save page  failed")
	}
	return c.Redirect(http.StatusSeeOther, "/"+p.Title)
}

func init() {
	gob.Register(&mdl.User{})
	setting.Initialize()
	ll, err := log.ParseLevel(setting.LogLevel)
	if err == nil {
		log.SetLevel(ll)
	}
}

func main() {
	mdl.SetupDb()

	store = session.NewCookieStore([]byte(setting.SessionSecret))

	e := echo.New()
	e.Use(session.Sessions("session", store))

	assetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())
	e.Get("/static/*", func(c echo.Context) error {
		http.StripPrefix("/static/", assetHandler).ServeHTTP(c.Response().(*standard.Response).ResponseWriter, c.Request().(*standard.Request).Request)
		return nil
	})
	e.Get("/favicon.ico", func(c echo.Context) error {
		http.StripPrefix("", assetHandler).ServeHTTP(c.Response().(*standard.Response).ResponseWriter, c.Request().(*standard.Request).Request)
		return nil
	})
	e.Use(setUser())
	if setting.ServerLogging {
		e.Use(serverLogger())
	}

	e.Get("/", root)
	e.Get("/*", wikiPage)
	e.Get("/s*", wikiPage)

	s := e.Group("/special")
	s.Get("/edit", edit)
	s.Post("/edit", savePage)
	s.Get("/history", history)
	s.Get("/recentchanges", recentChanges)
	s.Get("/pages", pages)
	s.Get("/users", users)
	s.Get("/register", register)
	s.Post("/register", registerHandle)
	s.Get("/login", login)
	s.Post("/login", loginHandle)
	s.Get("/logout", logout)
	s.Get("/random", random)
	s.Get("/delete", delete)
	s.Post("/delete", deleteHandle)
	p := s.Group("/preferences")
	p.Use(checkLoggedIn())
	p.Get("", prefs)
	p.Get("/password", prefsPasword)
	p.Post("/password", handlePrefsPassword)
	a := s.Group("/admin")
	a.Use(checkAdmin())
	a.Get("", admin)
	a.Post("", adminHandle)

	e.Use(middleware.Gzip())
	fmt.Println("Server started on port " + setting.HttpPort)
	e.Run(standard.New(":" + setting.HttpPort))
}
