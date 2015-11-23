package main

import (
	"net/http"
	"strings"
	"time"

	m "1bwiki/model"
	"1bwiki/tmpl"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mgutz/logxi/v1"
)

var logger log.Logger

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
	pv := m.GetPageView(n, t)

	if pv.NiceTitle != "" {
		return c.HTML(http.StatusOK, tmpl.Page(pv.NiceTitle, pv.Text))
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/B:edit?title="+n+t+"&action=edit")
}

func savePage(c *echo.Context) error {
	t := m.Text{Text: c.Form("text")}
	r := m.Revision{
		PageTitle: c.Form("title"),
		Comment:   c.Form("summary"),
		UserID:    1, // TODO :(
		UserText:  "pepp",
		Minor:     false,
		Deleted:   false,
		Len:       len(c.Form("text")),
		ParentID:  0,
		TimeStamp: time.Now().Unix(),
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
	logger = log.New("1bwiki")
	db, err := sqlx.Connect("sqlite3", "./1bwiki.db")
	if err != nil {
		panic(err)
	}
	// Convert to transaction
	db.Exec(`create table if not exists text (id integer primary KEY, text blob)`)
	db.Exec(`create table if not exists revision (id integer primary key,
			pagetitle text, textid integer, comment text, userid int,
			usertext text, minor integer, deleted integer, len integer,
			parentid integer, timestamp integer)`)
	db.Exec(`create table if not exists page (title text primary key,
			namespace text, nicetitle text, redirect integer, revisionid integer,
			len integer)`)
}

func main() {
	e := echo.New()
	e.Use(Logger())
	e.Use(fixURL())
	e.StripTrailingSlash()
	e.Static("/static", "static")

	e.Get("/", root)
	e.Get("/*", wikiPage)
	e.Post("/save", savePage)

	e.Get("/B:edit", edit)
	e.Get("/B:edits", edits)

	e.Run(":8000")
}
