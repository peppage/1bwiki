package main

import (
	"net/http"
	"strings"

	m "1bwiki/model"
	"1bwiki/tmpl"

	"github.com/labstack/echo"
)

func edit(c *echo.Context) error {
	n, t := parseTitle(c.Query("title"))
	if c.Query("action") == "edit" {
		pv := m.GetPageView(n, t)
		if pv.NiceTitle == "" {
			pv.NameSpace = n
			pv.Title = t
			pv.NiceTitle = strings.Replace(n+t, "_", " ", -1)
		}
		return c.HTML(http.StatusOK, tmpl.Editpage(pv))
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Not an acceptable action")
}

func edits(c *echo.Context) error {
	revs, err := m.GetRevisions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	return c.HTML(http.StatusOK, tmpl.Edits(revs))
}
