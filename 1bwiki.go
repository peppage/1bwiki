package main

import (
	"encoding/gob"
	"net/http"
	"strconv"
	"strings"
	"time"

	m "1bwiki/model"
	"1bwiki/tmpl"

	"github.com/gorilla/context"
	"github.com/labstack/echo"
	"github.com/mgutz/logxi/v1"
	"github.com/syntaqx/echo-middleware/session"
)

const secret = "Thisisatemporarysecret"

var logger log.Logger
var store = session.NewCookieStore([]byte(secret))

func cleanTitle(t string) string {
	f := string(t[0])
	t = strings.ToUpper(f) + string(t[1:])
	t = strings.Replace(t, "%20", "_", -1)
	t = strings.Replace(t, " ", "_", -1)
	return t
}

func parseTitle(title string) (string, string) {
	u := strings.Trim(title, "/")
	n := ""
	t := ""
	if strings.Contains(u, ":") {
		split := strings.Split(u, ":")
		n = split[0]
		t = split[1]
	} else {
		t = u
	}
	return n, t
}

func root(c *echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/Main_Page")
}

func wikiPage(c *echo.Context) error {
	n, t := parseTitle(c.Request().URL.String())

	ct := cleanTitle(t)

	if ct != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusMovedPermanently, "/"+n+ct)
	}

	pv := m.GetPageView(n, t)

	if pv.NiceTitle != "" {
		return c.HTML(http.StatusOK, tmpl.Page(pv.Title, pv.NiceTitle, pv.Html()))
	}
	if n != "" {
		n += ":"
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/special/action?title="+n+t+"&action=edit")
}

func savePage(c *echo.Context) error {
	t := m.Text{Text: c.Form("text")}
	l, err := strconv.Atoi(c.Form("len"))
	if err != nil {
		l = 0
		logger.Warn("save page len Atoi failed")
	}
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*m.User)
	if !ok {
		return logger.Error("User saving page is invalid", "user", u)
	}
	logger.Info("Saving page", "user", u)
	r := m.Revision{
		PageTitle: c.Form("title"),
		Comment:   c.Form("summary"),
		UserID:    u.ID,
		UserText:  u.Name,
		Minor:     false, // NOT WORKING
		Deleted:   false,
		Len:       len(c.Form("text")),
		ParentID:  0, // NOT WORKING
		TimeStamp: time.Now().Unix(),
		LenDiff:   len(c.Form("text")) - l,
	}
	if t.Verify() == nil && r.Verify() == nil {
		p := m.Page{
			Title:     c.Form("title"),
			Namespace: c.Form("namespace"),
			NiceTitle: strings.Replace(c.Form("title"), "_", " ", -1),
			Redirect:  false,
			Len:       len(c.Form("text")),
		}
		p.SavePage(t, r)
		return c.Redirect(http.StatusSeeOther, p.Title)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Save page not valid")
}

func init() {
	gob.Register(&m.User{})
	logger = log.New("1bwiki")
}

func main() {
	e := echo.New()
	e.Use(session.Sessions("session", store))
	e.Static("/static", "static")
	e.HTTP2(true)
	e.Use(setUser())
	e.Use(serverLogger())

	e.Get("/", root)
	e.Get("/*", wikiPage)
	e.Post("/save", savePage)

	s := e.Group("/special")
	s.Get("/action", action)
	s.Get("/edit", edit)
	s.Get("/history", history)
	s.Get("/recentchanges", recentChanges)
	s.Get("/register", register)
	s.Post("/register", registerHandle)
	s.Get("/login", login)
	s.Post("/login", loginHandle)
	s.Get("/logout", logout)
	s.Get("/prefences", prefs)

	http.ListenAndServe(":8000", context.ClearHandler(e))
}
