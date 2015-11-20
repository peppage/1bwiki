package main

import (
	"fmt"
	"net/http"
	"strings"

	m "1bwiki/model"
	"1bwiki/tmpl"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func wikiPage(c *echo.Context) error {
	fmt.Println(c.Request().URL)
	u := strings.Trim(c.Request().URL.String(), "/")
	n := ""
	t := ""
	if strings.Contains(u, ":") {
		split := strings.Split(u, ":")
		n = split[0]
		t = split[1]
	} else {
		t = u
	}
	p := m.GetPage(n, t)
	if p != nil {
		return c.String(http.StatusOK, "Page Exists")
	}
	return c.HTML(http.StatusOK, tmpl.Newpage(n, t))
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
		Len:       1,
		ParentID:  0,
		Sha1:      "aaaa",
	}
	if t.Verify() == nil && r.Verify() == nil {
		p := m.Page{
			Title:     c.Form("title"),
			Namespace: c.Form("namespace"),
			NiceTitle: strings.Replace(c.Form("title"), "_", " ", -1),
			Redirect:  false,
			Len:       1,
		}
		p.SavePage(t, r)
	}
	return c.String(http.StatusOK, "page not saved")
}

func init() {
	db, err := sqlx.Connect("sqlite3", "./1bwiki.db")
	if err != nil {
		panic(err)
	}
	db.Exec(`create table if not exists text (id integer primary KEY, text blob)`)
	db.Exec(`create table if not exists revision (id integer primary key,
			pagetitle text, textid integer, comment text, userid int,
			usertext text, minor integer, deleted integer, len integer,
			parentid integer, sha1 text)`)
	db.Exec(`create table if not exists page (title text primary key,
			namespace text, nicetitle text, redirect integer, revisionid integer,
			len integer)`)
}

func main() {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(fixURL())
	e.StripTrailingSlash()
	e.Static("/static", "static")

	e.Get("/*", wikiPage)
	e.Post("/save", savePage)

	e.Run(":8000")
}
