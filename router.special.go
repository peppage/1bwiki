package main

import (
	"net/http"
	"strings"

	m "1bwiki/model"
	"1bwiki/tmpl"
	"1bwiki/tmpl/special"

	"github.com/labstack/echo"
)

func action(c *echo.Context) error {
	n, t := parseTitle(c.Query("title"))
	ct := cleanTitle(t)
	if ct != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/special/action?title="+n+ct+"&action="+c.Query("action"))
	}

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
	} else if c.Query("oldid") != "" {
		pv := m.GetOldPageView(c.Query("oldid"))
		return c.HTML(http.StatusOK, special.Oldversion(pv))
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
